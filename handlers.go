package main

import (
	"encoding/json"
	"fmt"
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

	// Check if the request body is empty
	if r.Body == nil {
		http.Error(w, "Please send a request body", http.StatusBadRequest)
		return
	}

	// Check if the request body is too large
	if r.ContentLength > 1024 {
		http.Error(w, "Request body is too large", http.StatusRequestEntityTooLarge)
		return
	}

	// Decode the JSON request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println(req)

	// Make sure the request body is valid
	if req.UserId == "" {
		http.Error(w, "Missing userId", http.StatusBadRequest)
		return
	}

	// Make sure the request body is valid
	if req.Latitude < -90 || req.Latitude > 90 {
		http.Error(w, "Invalid latitude", http.StatusBadRequest)
		return
	}

	// Make sure the request body is valid
	if req.Longitude < -180 || req.Longitude > 180 {
		http.Error(w, "Invalid longitude", http.StatusBadRequest)
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
	// Respond with a JSON message
	json.NewEncoder(w).Encode(map[string]string{"message": "Location updated"})
}

func getLocationHandler(w http.ResponseWriter, r *http.Request) {
	// Get the userId from the query string
	userId := r.URL.Query().Get("userId")

	// Get the location from the store
	loc, ok := store.GetLocation(userId)
	if !ok {
		http.Error(w, "Location not found", http.StatusNotFound)
		return
	}

	// Marshal the struct to JSON
	jsonData, err := json.Marshal(loc)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
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
