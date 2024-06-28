package main

import (
	"encoding/json"
	"net/http"
	"pinflow-orbit/storage"
	"time"
)

func setLocationHandler(w http.ResponseWriter, r *http.Request) {
	// Parse JSON request body
	var req struct {
		UserId    string  `json:"userId"`
		Longitude float64 `json:"longitude"`
		Latitude  float64 `json:"latitude"`
	}

	// Decode the JSON request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get the current time
	now := time.Now().Unix()

	var loc storage.Location
	loc.Latitude = req.Latitude
	loc.Longitude = req.Longitude
	loc.LastUpdate = now

	// Update location in the store
	store.SetLocation(req.UserId, loc)

	w.WriteHeader(http.StatusOK)
}

func getAllLocationsHandler(w http.ResponseWriter, r *http.Request) {
	locations := store.GetAllLocations()

	// Marshal the struct to JSON
	jsonData, err := json.Marshal(locations)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
