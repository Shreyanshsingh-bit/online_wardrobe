package main
import(
	"encoding/json"
	"log"
	"fmt"
	"time"
	"net/http"
	"github.com/lib/pq"
)

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
func getRecommendationHandler(w http.ResponseWriter, r *http.Request){

	// fmt.Println("➡️ Step 1: Hit the /recommend endpoint!") // DEBUG PRINT

	if r.Method != http.MethodGet {
		http.Error(w, "Only Get method is allowed", http.StatusMethodNotAllowed)
		return
	}
	userIDStr := r.URL.Query().Get("user_id")
	lat := r.URL.Query().Get("lat")
	lon := r.URL.Query().Get("lon")

	// fmt.Printf("➡️ Step 2: Extracted Query Parameters -> User: %s, Lat: %s, Lon: %s\n", userIDStr, lat, lon) // DEBUG PRINT

	if userIDStr == "" || lat == "" || lon == ""{
		http.Error(w, "Missing required user_id, lat, and lon parameter", http.StatusBadRequest)
		return
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	weatherURL := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&current_weather=true", lat, lon)

	// fmt.Println("➡️ Step 3: Pinging External API:", weatherURL) // DEBUG PRINT

	resp, err := client.Get(weatherURL)
	if err != nil {
		log.Println("External weather API connection error:", err)
		http.Error(w, "Failed to fetch live weather telemetry", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close() // prevents network resource leaks

	// Decode the external API's response into our new OpenMeteoResponse struct

	var apiData OpenMeteoResponse
	err = json.NewDecoder(resp.Body).Decode(&apiData)
	if err != nil {
		log.Println("Error parsing external weather JSON:", err)
		http.Error(w, "Failed to parse live weather data", http.StatusInternalServerError)
		return
	}
	
	liveTemp := int(apiData.CurrentWeather.Temperature) // Converting float64 to int to match database
	liveIsRaining := apiData.CurrentWeather.WeatherCode >= 51 && apiData.CurrentWeather.WeatherCode <= 67

	// Fetches all clothes for this user from the database

	rows, err := db.Query(`SELECT id, user_id, image_url, category, sub_category, primary_color, material, min_temp_celsius, max_temp_celsius, is_waterproof, suitable_seasons, is_trending, created_at FROM clothing_items WHERE user_id = $1`, userIDStr)
	if err != nil {
		log.Println("Database query error:", err)
		http.Error(w, "Database failure", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
// here starts the recommendation algorithm(filtering loop)
// uses the live data
	var recommendedOutfit []ClothingItem
	
	for rows.Next() {
		var item ClothingItem
		err := rows.Scan(&item.ID, &item.UserID, &item.ImageURL, &item.Category, &item.SubCategory, &item.PrimaryColor, &item.Material, &item.MinTempCelsius, &item.MaxTempCelsius, &item.IsWaterproof, pq.Array(&item.SuitableSeasons), &item.IsTrending, &item.CreatedAt)
		if err != nil {
			log.Println("Row scan error:", err)
			http.Error(w, "Error reading clothing data", http.StatusInternalServerError)
			return
		}
		// the algo rules

		//matches live temperature
		if liveTemp < item.MinTempCelsius || liveTemp > item.MaxTempCelsius {
			continue 
		}

		// Matches live precipitation context
		if liveIsRaining && !item.IsWaterproof && (item.Category == "Outerwear" || item.Category == "Footwear") {
			continue 
		}

		// does it match the current season
		// seasonMatch := false
		// for _, s := range item.SuitableSeasons {
		// 	if s == weather.CurrentSeason {
		// 		seasonMatch = true
		// 		break
		// 	}
		// }
		// if !seasonMatch {
		// 	continue // Skip if not designed for this season
		// }
		// If it passed all filters, add it to the outfit!
		recommendedOutfit = append(recommendedOutfit, item)

		// returning the filtered outfit to the user
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(recommendedOutfit)
	}
}
// deleteClothingItemHandler removes a specific clothing item from the database by its ID
func deleteClothingItemHandler(w http.ResponseWriter, r *http.Request) {

	itemIDStr := r.URL.Query().Get("id")// extract items from url
	if itemIDStr == "" {
		http.Error(w, "Missing required id parameter", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(`DELETE FROM clothing_items WHERE id = $1`, itemIDStr) // the del statement
	if err != nil {
		log.Println("Database deletion error:", err)
		http.Error(w, "Database failure", http.StatusInternalServerError)
		return
	}

	//Verifying if the row actually existed
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error checking rows affected:", err)
		http.Error(w, "Database failure", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No clothing item found with that ID", http.StatusNotFound)
		return
	}

	//Responds with a successful message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Clothing item successfully removed from your closet"})
}


func filterOutfit(closet []ClothingItem, liveTemp int, isRaining bool) []ClothingItem {
	var recommended []ClothingItem

	for _, item := range closet {
		//Temperature Check
		if liveTemp < item.MinTempCelsius || liveTemp > item.MaxTempCelsius {
			continue
		}

		
		if isRaining && !item.IsWaterproof && (item.Category == "Outerwear" || item.Category == "Footwear") {
			continue
		}

		recommended = append(recommended, item)
	}

	return recommended
}