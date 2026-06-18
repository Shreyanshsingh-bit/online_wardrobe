package main

import (
	"database/sql"
	// "encoding/json"
	"fmt"
	"log"
	"net/http"

	// "github.com/lib/pq"   //imports postgres driver
)
// db is global variable for future APIs
var db *sql.DB


func main() {
	//The connection string uses database credentials
	connStr := "user=admin password=wardrobe123 dbname=wardrobe_db host=localhost sslmode=disable"

	//Open the connection
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open a database connection: ", err)
	}

	//Pinged the database to verify the connection is alive
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping the database: ", err)
	}

	fmt.Println("Successfully connected to the PostgreSQL Wardrobe Database!")

	// Registering the endpoint
	http.HandleFunc("/users", createdUserHandler) // tells server if someone goes to /users run this fn

	

	http.HandleFunc("/clothes", func(w http.ResponseWriter, r *http.Request){ // known as clothes routing block
		switch r.Method {
		case http.MethodGet:
			getClothingItemsHandler(w, r) // Directs to the Read fn
		
		case http.MethodPost:
			addClothingItemHandler(w, r) // directs to the create fn

		case http.MethodDelete: // hadles delete fn
			deleteClothingItemHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed) // if someone send a PUT, DELETE, or PATCH request, handle it here
		}

	})
	http.HandleFunc("/recommend", getRecommendationHandler) // recommendation doorway

	// Starts the server
	fmt.Println("Server is starting on port 8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server error: ", err)
	}
}