package main

import "time"

// User mirrors the 'users' table in your database
type User struct {
	ID                   int       `json:"id"`
	Name                 string    `json:"name"`
	Email                string    `json:"email"`
	PasswordHash         string    `json:"-"` // The "-" ensures the password is never sent to the frontend
	PreferredGenderStyle string    `json:"preferred_gender_style"`
	SizeTop              string    `json:"size_top"`
	SizeBottom           string    `json:"size_bottom"`
	CreatedAt            time.Time `json:"created_at"`
}

// ClothingItem mirrors the 'clothing_items' table in your database
type ClothingItem struct {
	ID               int       `json:"id"`
	UserID           int       `json:"user_id"`
	ImageURL         string    `json:"image_url"`
	Category         string    `json:"category"`
	SubCategory      string    `json:"sub_category"`
	PrimaryColor     string    `json:"primary_color"`
	Material         string    `json:"material"`
	MinTempCelsius   int       `json:"min_temp_celsius"`
	MaxTempCelsius   int       `json:"max_temp_celsius"`
	IsWaterproof     bool      `json:"is_waterproof"`
	SuitableSeasons  []string  `json:"suitable_seasons"` 
	IsTrending       bool      `json:"is_trending"`
	CreatedAt        time.Time `json:"created_at"`
}

// WeatherContext is a helper struct to handle incoming weather data from an API
type WeatherContext struct {
	TemperatureCelsius int    `json:"temperature_celsius"`
	IsRaining          bool   `json:"is_raining"`
	CurrentSeason      string `json:"current_season"`
}

