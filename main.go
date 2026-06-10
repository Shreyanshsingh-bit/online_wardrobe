package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/lib/pq"   //imports postgres driver
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
func addClothingItemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost{
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
	var newItem ClothingItem
	err := json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		http.Error(w, "Invalid JSON data provided", http.StatusBadRequest)
		return
	}
	sqlStatement := `
		INSERT INTO clothing_items 
		(user_id, image_url, category, sub_category, primary_color, material, min_temp_celsius, max_temp_celsius, is_waterproof, suitable_seasons, is_trending)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at`
	
	err = db.QueryRow(
		sqlStatement,
		newItem.UserID,
		newItem.ImageURL,
		newItem.Category,
		newItem.SubCategory,
		newItem.PrimaryColor,
		newItem.Material,
		newItem.MinTempCelsius,
		newItem.MaxTempCelsius,
		newItem.IsWaterproof,
		// Translates Go []string to SQL Array
		pq.Array(newItem.SuitableSeasons), 
		newItem.IsTrending,
	).Scan(&newItem.ID, &newItem.CreatedAt)

	if err != nil {
		log.Println("Database insert error:", err)
		http.Error(w, "Failed to add clothing item. Ensure that user ID exists.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newItem)

	
}
func getClothingItemsHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet{
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}
	userIDStr := r.URL.Query().Get("user_id") //Read the user_id from the URL query parameters
	if userIDStr == ""{
		http.Error(w, "Missing required user_id parameter", http.StatusBadRequest)
		return

	}
	// query the database for all items belonging to this user
	rows, err := db.Query(`
		SELECT id, user_id, image_url, category, sub_category, primary_color, material, min_temp_celsius, max_temp_celsius, is_waterproof, suitable_seasons, is_trending, created_at 
		FROM clothing_items 
		WHERE user_id = $1`, userIDStr)

	if err != nil {
		log.Println("Database query error:", err)
		http.Error(w, "Database failure", http.StatusInternalServerError)
		return
	}

	defer rows.Close() // clears the database connection if already fn esists

	items := []ClothingItem{} // returns [] instead of null
	for rows.Next() { // loops the conveyor belt of database rows
		var item ClothingItem

		err:= rows.Scan(
			&item.ID,
			&item.UserID,
			&item.ImageURL,
			&item.Category,
			&item.SubCategory,
			&item.PrimaryColor,
			&item.Material,
			&item.MinTempCelsius,
			&item.MaxTempCelsius,
			&item.IsWaterproof,
			pq.Array(&item.SuitableSeasons),// using pq.array to unpack the SQL array
			&item.IsTrending,
			&item.CreatedAt,
		)
		if err != nil {
			log.Println("Row scan error:", err)
			http.Error(w, "Error parsing wardrobe data", http.StatusInternalServerError)
			return
		}
		// append the freshly populated item to our slice
		items = append(items, item)


	}	
	// ensures the loop didnt stop early because of hidden network issue
	if err = rows.Err(); err != nil {
		log.Println("Rows iteration error:", err)
		http.Error(w, "Error processing wardrobe data", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json") // sends complete array as a clean json response
	w.WriteHeader(http.StatusOK) // Status code 200 OK
	json.NewEncoder(w).Encode(items)
}

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

	// Registering the new clothes doorway
	// http.HandleFunc("/clothes", addClothingItemHandler)

	// Registering the endpoints to fetch clothes
	http.HandleFunc("/clothes", func(w http.ResponseWriter, r *http.Request){ // since we use the same url path for both adding and viewing
		if r.Method == http.MethodGet { // routing the traffic based on the http method
			getClothingItemsHandler(w, r)
		} else if r.Method == http.MethodPost {
			addClothingItemHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Starts the server
	fmt.Println("Server is starting on port 8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server error: ", err)
	}
}