package main

import (
	"log"
	"net/http"
	"packcalculator/endpoints"
	"packcalculator/helpers"
	"packcalculator/logger"
	"packcalculator/storage"

	"github.com/joho/godotenv"
)

func main() {
	// Get env variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Setup logger
	logger.SetupLogger()

	// Setup redis client
	err = storage.SetupRedisStorageClient()
	if err != nil {
		logger.Instance.Error("Failed to setup redis storage client", err)
	}

	// Setup api endpoints
	http.HandleFunc("/packs", endpoints.HandlePacksEndpoints)
	http.HandleFunc("/calculate", endpoints.HandleCalculateEndpoints)

	// Start application
	logger.Instance.Info("Starting API")
	err = http.ListenAndServe(":8080", helpers.WithCORS(http.DefaultServeMux))
	if err != nil {
		logger.Instance.Fatal("Failed to start server", err)
	}
}
