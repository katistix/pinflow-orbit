package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"sync"
	"time"
)

// Location represents a geographical location with latitude and longitude.
type Location struct {
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	LastUpdate int64   `json:"last_update"`
}

// LocationStore manages in-memory storage of Location objects.
type LocationStore struct {
	lock      sync.RWMutex
	locations map[string]Location
}

// NewLocationStore creates a new instance of LocationStore.
func NewLocationStore() *LocationStore {
	return &LocationStore{
		locations: make(map[string]Location),
	}
}

// OPERATIONS

// SetLocation updates or sets the location for a given key.
func (s *LocationStore) SetLocation(key string, loc Location) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.locations[key] = loc
}

// GetLocation retrieves the location for a given key.
func (s *LocationStore) GetLocation(key string) (Location, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	loc, ok := s.locations[key]
	return loc, ok
}

// DeleteLocation removes a location entry for a given key.
func (s *LocationStore) DeleteLocation(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.locations, key)
}

// GetAllLocations retrieves all locations in the store.
func (s *LocationStore) GetAllLocations() map[string]Location {
	s.lock.RLock()
	defer s.lock.RUnlock()
	// Create a copy of the map to avoid concurrent map read/write issues
	locations := make(map[string]Location, len(s.locations))
	for k, v := range s.locations {
		locations[k] = v
	}
	return locations
}

func AuthorizeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Check if the Authorization header is set
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Check if the Authorization header is the same as the value in the .env file
		if r.Header.Get("Authorization") != "Bearer "+os.Getenv("API_KEY") {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Hello, World!"}`))

}

var store *LocationStore

func main() {

	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Create a new ServeMux
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	mux.HandleFunc("/protected", AuthorizeMiddleware(http.HandlerFunc(handler)).ServeHTTP)
	// Create a new LocationStore
	store := NewLocationStore()

	// ROUTES
	mux.HandleFunc("GET /locations", func(w http.ResponseWriter, r *http.Request) {
		locations := store.GetAllLocations()

		// Marshal the struct to JSON
		jsonData, err := json.Marshal(locations)
		if err != nil {
			http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})

	mux.HandleFunc("POST /set", func(w http.ResponseWriter, r *http.Request) {
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

		var loc Location
		loc.Latitude = req.Latitude
		loc.Longitude = req.Longitude
		loc.LastUpdate = now

		// Update location in the store
		store.SetLocation(req.UserId, loc)

		w.WriteHeader(http.StatusOK)
	})

	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", mux)
}
