package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"

	"techtarget_project/models"
	"techtarget_project/vault"
)

type ProvisionHandler struct {
	DB           *sqlx.DB
	VaultClient  *vault.Client
	UniversalKey string
}

func (h *ProvisionHandler) Provision(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode request body
	var req models.ProvisionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate required fields
	if req.SerialNumber == "" {
		http.Error(w, "serial_number is required", http.StatusBadRequest)
		return
	}
	if req.Model == "" {
		http.Error(w, "model is required", http.StatusBadRequest)
		return
	}
	if req.Firmware == "" {
		http.Error(w, "firmware is required", http.StatusBadRequest)
		return
	}
	if req.UniversalKey == "" {
		http.Error(w, "universal_key is required", http.StatusBadRequest)
		return
	}

	// Authenticate universal key
	if req.UniversalKey != h.UniversalKey {
		http.Error(w, "Unauthorized: invalid universal key", http.StatusUnauthorized)
		return
	}

	// Check if device already provisioned
	var exists bool
	err := h.DB.Get(&exists, "SELECT EXISTS(SELECT 1 FROM device_info WHERE serial_number = $1)", req.SerialNumber)
	if err != nil {
		log.Printf("DB error checking device serial_number=%s: %v", req.SerialNumber, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Device already provisioned", http.StatusConflict)
		return
	}

	// Generate unique API key
	apiKey, err := vault.GenerateAPIKey()
	if err != nil {
		log.Printf("Failed to generate API key for serial_number=%s: %v", req.SerialNumber, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Insert new device into device_info with the generated API key
	_, err = h.DB.Exec(`
		INSERT INTO device_info (serial_number, model, firmware, api_key, last_seen)
		VALUES ($1, $2, $3, $4, NOW())
	`, req.SerialNumber, req.Model, req.Firmware, apiKey)
	if err != nil {
		log.Printf("DB error inserting device serial_number=%s: %v", req.SerialNumber, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp := models.ProvisionResponse{
		Status: true,
		APIKey: apiKey,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode provision response: %v", err)
	}
}