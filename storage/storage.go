package storage

import "sync"

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
