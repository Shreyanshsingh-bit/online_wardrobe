package main

import (
	"testing"
)
func TestFilterOutfit(t *testing.T) {
	// 1. Setup a Fake Closet in memory (No database needed!)
	mockCloset := []ClothingItem{
		{ID: 1, Category: "Outerwear", MinTempCelsius: -10, MaxTempCelsius: 15, IsWaterproof: true},  // Winter Jacket
		{ID: 2, Category: "Tops",      MinTempCelsius: 20,  MaxTempCelsius: 45, IsWaterproof: false}, // Summer T-Shirt
		{ID: 3, Category: "Outerwear", MinTempCelsius: -10, MaxTempCelsius: 15, IsWaterproof: false}, // Wool Coat (Not waterproof)
	}

	rainyResult := filterOutfit(mockCloset, 8, true)

	if len(rainyResult) != 1 {
		t.Fatalf("Scenario A Failed: Expected 1 item, got %d", len(rainyResult))
	}
	if rainyResult[0].ID != 1 {
		t.Errorf("Scenario A Failed: Expected Winter Jacket (ID 1), got ID %d", rainyResult[0].ID)
	}

	
	hotResult := filterOutfit(mockCloset, 35, false)

	if len(hotResult) != 1 {
		t.Fatalf("Scenario B Failed: Expected 1 item, got %d", len(hotResult))
	}
	if hotResult[0].ID != 2 {
		t.Errorf("Scenario B Failed: Expected Summer T-Shirt (ID 2), got ID %d", hotResult[0].ID)
	}
}