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
	db *sql.DB
)

// Customer - struct for customer data
type Customer struct {
	ID        int    `json:"Id"`
	Surname   string `json:"Title"`
	Givenname string `json:"desc"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the API!")
	fmt.Println("Endpoint Hit: Main API page")
}

func returnCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID, err := strconv.ParseInt(vars["id"], 10, 32)

	if err != nil {
		println("Error while parsing a customer id!")
	} else {
		// Prepare statement for reading data
		stmtOut, err := db.Prepare("SELECT Customers_ID, Surname, Givenname FROM Customers WHERE Customer_ID = ?")
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		defer stmtOut.Close()

		var customerData Customer // we "scan" the result in here

		// Query the square-number of the customer id store it in customerdata
		err = stmtOut.QueryRow(customerID).Scan(&customerData.ID, &customerData.Surname, customerData.Givenname)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		fmt.Printf("The name of customer %d is: %s %s", customerData.ID, customerData.Givenname, customerData.Surname)

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
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {

	db, err := sql.Open("mysql", "admin:admin@tcp(123.4.5.6:3306)/database-1")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	handleRequests()
}
