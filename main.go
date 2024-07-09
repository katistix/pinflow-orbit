package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"pinflow-orbit/storage"
)

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

var store *storage.LocationStore

func main() {

	// Only if not in production, load the .env file
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			fmt.Println("Error loading .env file")
			os.Exit(1)
		}
	}

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Create a new LocationStore
	store = storage.NewLocationStore()

	// HEALTH CHECK
	mux.Handle("GET /health", http.HandlerFunc(healthCheckHandler))

	// ROUTES
	mux.Handle("GET /locations", AuthorizeMiddleware(http.HandlerFunc(getAllLocationsHandler)))
	mux.Handle("GET /location", AuthorizeMiddleware(http.HandlerFunc(getLocationHandler)))
	mux.Handle("POST /set", AuthorizeMiddleware(http.HandlerFunc(setLocationHandler)))
	mux.Handle("POST /delete", AuthorizeMiddleware(http.HandlerFunc(deleteLocationHandler)))

	// Serve the API
	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", mux)
}
