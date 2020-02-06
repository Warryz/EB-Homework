// main.go
package main

// See: https://github.com/TutorialEdge/create-rest-api-in-go-tutorial

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var (
	db  *sql.DB
	err error
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
	MeasureValue string `json:"MeasureValue"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the API!")
	fmt.Println("Endpoint Hit: Main API page")
}

func returnCustomerData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID, parseErr := strconv.ParseInt(vars["id"], 10, 32)

	if parseErr != nil {
		println("dbError while parsing a customer id!")
	} else {
		// Prepare statement for reading data
		stmtOut, dbErr := db.Prepare("SELECT Measure_ID, Measure_Date, Value FROM Readings WHERE Customers_ID_FK = ?;")
		if dbErr != nil {
			fmt.Println("Error while creating the sql statement")
		}
		defer stmtOut.Close()

		var customerData Customer // we "scan" the result in here

		// Query the customer id store it in customerdata

		// Continue here!!! Need to save the data from the db with a loop
		dbErr = stmtOut.QueryRow(customerID).Scan(&customerData.ID, &customerData.Surname, &customerData.Givenname)
		if dbErr != nil {
			fmt.Println("unable to query user", customerID, dbErr)
		} else {
			fmt.Printf("The name of customer %d is: %s %s", customerData.ID, customerData.Givenname, customerData.Surname)

			json.NewEncoder(w).Encode(customerData)
		}
	}

}

func returnCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID, parseErr := strconv.ParseInt(vars["id"], 10, 32)

	if parseErr != nil {
		println("dbError while parsing a customer id!")
	} else {
		// Prepare statement for reading data
		stmtOut, dbErr := db.Prepare("SELECT Customers_ID, Surname, Givenname FROM Customers WHERE Customers_ID = ?")
		if dbErr != nil {
			fmt.Println("Error while creating the sql statement")
		}
		defer stmtOut.Close()

		var customerData Customer // we "scan" the result in here

		// Query the customer id store it in customerdata
		dbErr = stmtOut.QueryRow(customerID).Scan(&customerData.ID, &customerData.Surname, &customerData.Givenname)
		if dbErr != nil {
			fmt.Println("unable to query user", customerID, dbErr)
		} else {
			fmt.Printf("The name of customer %d is: %s %s", customerData.ID, customerData.Givenname, customerData.Surname)

			json.NewEncoder(w).Encode(customerData)
		}
	}

}

func createNewCustomer(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// unmarshal this into a new Article struct
	// append this to our Articles array.
	reqBody, _ := ioutil.ReadAll(r.Body)
	var cus Customer
	json.Unmarshal(reqBody, &cus)
	// update our global Articles array to include
	// our new Article

	json.NewEncoder(w).Encode(cus)
}

// func deleteArticle(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id := vars["id"]

// 	for index, article := range Articles {
// 		if article.ID == id {
// 			Articles = append(Articles[:index], Articles[index+1:]...)
// 		}
// 	}

// }

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	// myRouter.HandleFunc("/articles", returnAllArticles)
	// myRouter.HandleFunc("/article", createNewArticle).Methods("POST")
	// myRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
	myRouter.HandleFunc("/customer/{id}", returnCustomer)
	myRouter.HandleFunc("/customerdata/{id}", returnCustomerData)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {

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
	handleRequests()
}
