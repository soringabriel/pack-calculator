package endpoints

import (
	"encoding/json"
	"net/http"
	"packcalculator/storage"
)

// Request and response structs
type PostPackRequest struct {
	Packs []int `json:"packs"`
}
type PostPackResponse struct {
	Success bool `json:"success"`
}
type GetPackResponse struct {
	Packs []int `json:"packs"`
}

// Handle requests on /packs
func HandlePacksEndpoints(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetPacksEndpoint(w, r)
	case http.MethodPost:
		PostPacksEndpoints(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "Method not allowed"}`))
		return
	}
}

// Handle GET requests on /packs
// Returns a list of pack sizes from redis
func GetPacksEndpoint(w http.ResponseWriter, r *http.Request) {
	// Get packs from redis
	packs, err := storage.RedisStorageClient.GetPackSizes()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to get pack sizes"}`))
		return
	}

	// Build response
	result := GetPackResponse{Packs: packs}
	jsonResult, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to marshal pack sizes"}`))
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonResult))
}

// Handle POST requests on /packs
// Sets pack sizes in redis
func PostPacksEndpoints(w http.ResponseWriter, r *http.Request) {
	// Get request struct
	var postPackRequest PostPackRequest
	err := json.NewDecoder(r.Body).Decode(&postPackRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Failed to decode request body"}`))
		return
	}

	// Validate request
	if len(postPackRequest.Packs) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Packs is required"}`))
		return
	}
	for _, pack := range postPackRequest.Packs {
		if pack <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Pack sizes must be greater than 0"}`))
			return
		}
	}

	// Set packs in redis
	err = storage.RedisStorageClient.SetPackSizes(r.Context(), postPackRequest.Packs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to set pack sizes"}`))
		return
	}

	// Build response
	result := PostPackResponse{Success: true}
	jsonResult, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to marshal response"}`))
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonResult))
}
