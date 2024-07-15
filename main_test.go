package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"pinflow-orbit/storage"
	"testing"
	"time"
)

func init() {
	// Load environment variables for testing
	err := os.Setenv("API_KEY", "testkey")
	if err != nil {
		return
	}

	store = storage.NewLocationStore()
}

func TestAuthorizeMiddleware(t *testing.T) {
	handler := AuthorizeMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	tests := []struct {
		name       string
		apiKey     string
		wantStatus int
	}{
		{"No API Key", "", http.StatusForbidden},
		{"Invalid API Key", "Bearer invalidkey", http.StatusForbidden},
		{"Valid API Key", "Bearer testkey", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			if tt.apiKey != "" {
				req.Header.Set("Authorization", tt.apiKey)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantStatus)
			}
		})
	}
}

func TestHealthCheckHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthCheckHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestSetLocationHandler(t *testing.T) {
	tests := []struct {
		name       string
		body       interface{}
		wantStatus int
		authHeader string
	}{
		{
			name:       "Valid Request",
			body:       map[string]interface{}{"userId": "user1", "longitude": 10.0, "latitude": 20.0},
			wantStatus: http.StatusOK,
			authHeader: "Bearer testkey",
		},
		{
			name:       "Empty Body",
			body:       nil,
			wantStatus: http.StatusBadRequest,
			authHeader: "Bearer testkey",
		},
		{
			name:       "Invalid JSON",
			body:       "invalid json",
			wantStatus: http.StatusBadRequest,
			authHeader: "Bearer testkey",
		},
		{
			name:       "Missing userId",
			body:       map[string]interface{}{"longitude": 10.0, "latitude": 20.0},
			wantStatus: http.StatusBadRequest,
			authHeader: "Bearer testkey",
		},
		{
			name:       "Invalid Latitude",
			body:       map[string]interface{}{"userId": "user1", "longitude": 10.0, "latitude": 200.0},
			wantStatus: http.StatusBadRequest,
			authHeader: "Bearer testkey",
		},
		{
			name:       "Invalid Longitude",
			body:       map[string]interface{}{"userId": "user1", "longitude": 200.0, "latitude": 20.0},
			wantStatus: http.StatusBadRequest,
			authHeader: "Bearer testkey",
		},
		{
			name:       "Missing Authorization Header",
			body:       map[string]interface{}{"userId": "user1", "longitude": 10.0, "latitude": 20.0},
			wantStatus: http.StatusForbidden,
			authHeader: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			if tt.body != nil {
				body, _ = json.Marshal(tt.body)
			}

			req, err := http.NewRequest("POST", "/set", bytes.NewReader(body))
			if err != nil {
				t.Fatal(err)
			}
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			rr := httptest.NewRecorder()
			handler := AuthorizeMiddleware(http.HandlerFunc(setLocationHandler))

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantStatus)
			}
		})
	}
}

func TestGetLocationHandler(t *testing.T) {
	// Setup: add a location for testing
	store.SetLocation("user1", storage.Location{Latitude: 20.0, Longitude: 10.0, LastUpdate: time.Now().Unix()})

	tests := []struct {
		name       string
		userId     string
		wantStatus int
	}{
		{"Location Found", "user1", http.StatusOK},
		{"Location Not Found", "user2", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/location?userId="+tt.userId, nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Authorization", "Bearer testkey")

			rr := httptest.NewRecorder()
			handler := AuthorizeMiddleware(http.HandlerFunc(getLocationHandler))

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantStatus)
			}
		})
	}
}

func TestGetAllLocationsHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/locations", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer testkey")

	rr := httptest.NewRecorder()
	handler := AuthorizeMiddleware(http.HandlerFunc(getAllLocationsHandler))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Additional checks for response body can be added here if needed
}

func TestDeleteLocationHandler(t *testing.T) {
	// Setup: add a location for testing
	store.SetLocation("user1", storage.Location{Latitude: 20.0, Longitude: 10.0, LastUpdate: time.Now().Unix()})

	tests := []struct {
		name       string
		body       interface{}
		wantStatus int
		authHeader string
	}{
		{
			name:       "Valid Request",
			body:       map[string]interface{}{"userId": "user1"},
			wantStatus: http.StatusOK,
			authHeader: "Bearer testkey",
		},
		{
			name:       "Empty Body",
			body:       nil,
			wantStatus: http.StatusBadRequest,
			authHeader: "Bearer testkey",
		},
		{
			name:       "Invalid JSON",
			body:       "invalid json",
			wantStatus: http.StatusBadRequest,
			authHeader: "Bearer testkey",
		},
		{
			name:       "Missing userId",
			body:       map[string]interface{}{},
			wantStatus: http.StatusBadRequest,
			authHeader: "Bearer testkey",
		},
		{
			name:       "Missing Authorization Header",
			body:       map[string]interface{}{"userId": "user1"},
			wantStatus: http.StatusForbidden,
			authHeader: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			if tt.body != nil {
				body, _ = json.Marshal(tt.body)
			}

			req, err := http.NewRequest("POST", "/delete", bytes.NewReader(body))
			if err != nil {
				t.Fatal(err)
			}
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			rr := httptest.NewRecorder()
			handler := AuthorizeMiddleware(http.HandlerFunc(deleteLocationHandler))

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantStatus)
			}
		})
	}
}
