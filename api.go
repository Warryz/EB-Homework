package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bradfitz/gomemcache/memcache"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
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
		// Try reading the data from the memcached server
		key := fmt.Sprintf("customerReadings_id_%d", customerID)

		// mc.Get(&memcache.Item{Key: key, Value: []byte(b)})
		it, memErr := mc.Get(key)
		if memErr != nil {
			fmt.Printf("No data for customer id %d in memcached: %s", customerID, memErr)

			// Prepare statement for reading data
			stmtOut, dbErr := db.Prepare("SELECT Measure_ID, Measure_Date, Value FROM Readings WHERE Customers_ID_FK = ?;")
			if dbErr != nil {
				fmt.Println("Error while creating the sql statement")
			}
			defer stmtOut.Close()

			// Query the customer id store it in customerdata

			rows, dbErr := stmtOut.Query(customerID)
			defer rows.Close()

			customReadingsList := myReadings{}
			var customerReadings Readings

			if dbErr != nil {
				fmt.Println("unable to query user data", customerID, dbErr)
			} else {
				for rows.Next() {

					err := rows.Scan(&customerReadings.MeasureID, &customerReadings.MeasureDate, &customerReadings.MeasureValue)
					if err != nil {
						log.Fatal(err)
					}
					// https://stackoverflow.com/questions/18042439/go-append-to-slice-in-struct
					customReadingsList.AddItem(customerReadings)
				}
				json.NewEncoder(w).Encode(customReadingsList)
			}
		} else {
			// Output the memcached data.
			w.Write(it.Value)
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

func handleRequests(webserverPort int) {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/customer/{id}", returnCustomer)
	myRouter.HandleFunc("/customerdata/{id}", returnCustomerData)

	// Needed to disable connection timeouts
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", webserverPort),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      myRouter,
	}

	srv.SetKeepAlivesEnabled(false)
	fmt.Printf("Starting webserver on port: %d", webserverPort)
	log.Fatal(srv.ListenAndServe())
}

func loadCustomerDataIntoMemory() {
	// Load the customer data into cached
	stmtOut, dbErr := db.Prepare("SELECT * FROM Customers;")
	if dbErr != nil {
		fmt.Println("Error while creating the sql statement")
	}
	defer stmtOut.Close()

	// Query the customer id store it in customerdata
	rows, dbErr := stmtOut.Query()
	defer rows.Close()

	var customerData Customer // we "scan" the result in here

	if dbErr != nil {
		fmt.Println("unable to load user into memcached", dbErr)
	} else {

		for rows.Next() {
			// ID, Surname, givenname
			err := rows.Scan(&customerData.ID, &customerData.Surname, &customerData.Givenname)
			if err != nil {
				fmt.Println("unable to parse user row into memcached", err)
			}

			// Save the loaded data to memcached by converting it to json: https://stackoverflow.com/questions/8270816/converting-go-struct-to-json
			b, err := json.Marshal(customerData)
			if err != nil {
				fmt.Println(err)
				continue
			}
			// Format the key and
			key := fmt.Sprintf("customerData_id_%d", customerData.ID)
			//  Save the data to memcached servers
			mc.Set(&memcache.Item{Key: key, Value: []byte(b)})
		}
		fmt.Println("Finished cache creation for customer data.")
	}

}

func loadReadingsDataIntoMemory() {
	// Load the readings data into cache

	// Does not work, we need this user by user

	stmtOut, dbErr := db.Prepare("SELECT DISTINCT Customers_ID_FK FROM Readings;")
	if dbErr != nil {
		fmt.Println("Error while creating the sql statement")
	}
	defer stmtOut.Close()

	// Query the customer id store it in customerdata
	rows, dbErr := stmtOut.Query()
	defer rows.Close()

	if dbErr != nil {
		fmt.Println("unable to query user ids for caching memcached", dbErr)
	} else {
		// Load the user ids
		for rows.Next() {
			customReadingsList := myReadings{}
			var customerReadings Readings
			var customerID int

			err := rows.Scan(&customerID)
			if err != nil {
				// log.Fatal(err)
				fmt.Println("Could not parse user id")
				continue
			}

			// stmtOut, dbErr := db.Prepare("SELECT DISTINCT Customers_ID_FK FROM Readings;")
			stmtOut, dbErr := db.Prepare("SELECT Measure_ID, Measure_Date, Value FROM Readings where Customers_ID_FK = ?;")
			if dbErr != nil {
				fmt.Println("Error while creating the sql statement")
			}
			defer stmtOut.Close()

			// Query the customer id store it in customerdata
			readingRows, dbErr := stmtOut.Query(customerID)
			defer rows.Close()

			if dbErr != nil {
				fmt.Println("unable to query user ids for caching memcached", dbErr)
			} else {
				for readingRows.Next() {
					// IDFK, MeasureID, Date, Value
					err := readingRows.Scan(&customerReadings.MeasureID, &customerReadings.MeasureDate, &customerReadings.MeasureValue)
					if err != nil {
						log.Fatal(err)
					}
					// https://stackoverflow.com/questions/18042439/go-append-to-slice-in-struct
					customReadingsList.AddItem(customerReadings)
				}
				// Save the loaded data to memcached by converting it to json: https://stackoverflow.com/questions/8270816/converting-go-struct-to-json
				b, err := json.Marshal(customReadingsList)
				if err != nil {
					fmt.Println(err)
					continue
				}
				// Format the key and
				key := fmt.Sprintf("customerReadings_id_%d", customerID)
				//  Save the data to memcached servers
				mc.Set(&memcache.Item{Key: key, Value: []byte(b)})
			}
		}
	}
	fmt.Println("Finished cache creation for customer readings.")
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

	// Load all from the db into memcached
	loadCustomerDataIntoMemory()
	loadReadingsDataIntoMemory()

	// Print when ready to serve
	fmt.Println("Ready to serve traffic...")
}

func main() {
	// Port flag
	portPtr := flag.Int("port", 10000, "Port for the webserver to start")
	flag.Parse()

	initSetup()
	handleRequests(*portPtr)
}
