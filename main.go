package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"   //imports postgres driver
)
// db is global variable for future APIs
var db *sql.DB
func createdUserHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid JSON data provided", http.StatusBadRequest)
		return
	}
	sqlStatement := `
		INSERT INTO users (name, email, password_hash, preferred_gender_style, size_top, size_bottom)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`

	err = db.QueryRow(
		sqlStatement, 
		newUser.Name, 
		newUser.Email, 
		newUser.PasswordHash, 
		newUser.PreferredGenderStyle, 
		newUser.SizeTop, 
		newUser.SizeBottom,
	).Scan(&newUser.ID, &newUser.CreatedAt)

	if err != nil {
		log.Println("Database insert error:", err)
		http.Error(w, "Failed to create user. The email might already exist.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // This sends the standard "201 Created" success code
	json.NewEncoder(w).Encode(newUser)
}

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

	// Registering the endpoint
	http.HandleFunc("/users", createdUserHandler) // tells server if someone goes to /users run this fn
	
	// Starts the server
	fmt.Println("Server is starting on port 8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server error: ", err)
	}
}