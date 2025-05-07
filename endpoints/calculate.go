package endpoints

import (
	"encoding/json"
	"net/http"
	"packcalculator/helpers"
	"packcalculator/storage"

	"github.com/gorilla/schema"
)

// Request and response structs
type GetCalculateRequest struct {
	Amount *int `json:"amount"`
}
type GetCalculateResponse struct {
	Amount int         `json:"amount"`
	Result map[int]int `json:"result"`
}

// Handle requests on /calculate
func HandleCalculateEndpoints(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetCalculateEndpoint(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "Method not allowed"}`))
		return
	}
}

// Handle GET requests on /calculate
// Returns the number of packs of each size that can be used to make amount
func GetCalculateEndpoint(w http.ResponseWriter, r *http.Request) {
	// Get request struct
	var getCalculateRequest GetCalculateRequest
	decoder := schema.NewDecoder()
	decoder.Decode(&getCalculateRequest, r.URL.Query())

	// Validate request
	if getCalculateRequest.Amount == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Amount is required"}`))
		return
	}
	if *getCalculateRequest.Amount <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Amount must be greater than 0"}`))
		return
	}

	// Get packs from redis
	packs, err := storage.RedisStorageClient.GetPackSizes()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to get pack sizes"}`))
		return
	}
	if len(packs) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "No pack sizes found"}`))
		return
	}

	// Calculate
	result := helpers.FindOptimalPackCombination(packs, *getCalculateRequest.Amount)

	// Build response
	response := GetCalculateResponse{
		Amount: *getCalculateRequest.Amount,
		Result: result,
	}
	jsonResult, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to marshal response"}`))
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonResult))
}
