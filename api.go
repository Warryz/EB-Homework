package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bradfitz/gomemcache/memcache"
	_ "github.com/go-sql-driver/mysql"
)

var (
	db  *sql.DB
	err error

	// Memcached variable
	mc = *memcache.New("hausarbeit-eb-memcached.dldis0.cfg.euc1.cache.amazonaws.com:11211")
)

// Customer - struct for customer data
type Customer struct {
	ID        int    `json:"Id"`
	Surname   string `json:"Surname"`
	Givenname string `json:"Givenname"`
}

// Readings - struct for read data
type Readings struct {
	MeasureID    int    `json:"MeasureID"`
	MeasureDate  string `json:"MeasureDate"`
	MeasureValue int    `json:"MeasureValue"`
}

type myReadings struct {
	Measures []Readings
}

func (reading *myReadings) AddItem(item Readings) {
	reading.Measures = append(reading.Measures, item)
}

// Similarly add a helper for send responses relating to client errors.
func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func returnCustomerData(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// ID may only contain numbers
	var idRegExp = regexp.MustCompile(`[0-9]`)

	// Parse the id from the query string
	ID := req.QueryStringParameters["id"]

	// Check if the provided ID is valid
	if !idRegExp.MatchString(ID) {
		return clientError(http.StatusBadRequest)
	}
	key := fmt.Sprintf("customerReadings_id_%s", ID)

	it, memErr := mc.Get(key)
	if memErr != nil {
		fmt.Printf("No data for customer id %s in memcached: %s", ID, memErr)

		// Create a db connection
		initSetup()

		// Prepare statement for reading data
		stmtOut, dbErr := db.Prepare("SELECT Measure_ID, Measure_Date, Value FROM Readings WHERE Customers_ID_FK = ?;")
		if dbErr != nil {
			fmt.Println("Error while creating the sql statement")
		}
		defer stmtOut.Close()

		// Query the customer id store it in customerdata
		rows, dbErr := stmtOut.Query(ID)
		defer rows.Close()

		customReadingsList := myReadings{}
		var customerReadings Readings

		if dbErr != nil {
			fmt.Println("unable to query user data", ID, dbErr)
		} else {
			for rows.Next() {

				err := rows.Scan(&customerReadings.MeasureID, &customerReadings.MeasureDate, &customerReadings.MeasureValue)
				if err != nil {
					log.Fatal(err)
				}
				// https://stackoverflow.com/questions/18042439/go-append-to-slice-in-struct
				customReadingsList.AddItem(customerReadings)
			}
			json, err := json.Marshal(customReadingsList)
			if err != nil {
				return events.APIGatewayProxyResponse{
					StatusCode: http.StatusOK,
					Body:       string(json),
				}, nil
			}
		}
	}

	// Return the events and a http 200 code.
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(it.Value),
	}, nil

}

func initSetup() {

	db, err = sql.Open("mysql", "admin:admin@tcp(123.4.5.6)/hausarbeit")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	} else {
		fmt.Println("DB connection established!")
	}

}

func main() {
	// Start the Lambda Handler
	lambda.Start(returnCustomerData)
}
