package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"   //imports postgres driver
)
// db is global variable for future APIs
var db *sql.DB

func main() {
	//The connection string using your database credentials
	connStr := "user=admin password=wardrobe123 dbname=wardrobe_db host=localhost sslmode=disable"

	//Open the connection
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open a database connection: ", err)
	}

	//Ping the database to verify the connection is alive
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping the database: ", err)
	}

	fmt.Println("Successfully connected to the PostgreSQL Wardrobe Database!")

	// Starts the server
	fmt.Println("Server is starting on port 8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server error: ", err)
	}
}