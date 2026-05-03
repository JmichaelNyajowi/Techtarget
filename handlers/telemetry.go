package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"

	"techtarget_project/models"
)

type TelemetryHandler struct {
	DB *sqlx.DB
}

func (h *TelemetryHandler) Submit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.TelemetryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate required fields
	if req.Info.Serial == "" {
		http.Error(w, "info.serial is required", http.StatusBadRequest)
		return
	}
	if req.Info.APIKey == "" {
		http.Error(w, "info.api_key is required", http.StatusBadRequest)
		return
	}

	// Fetch stored API key for this device
	var storedAPIKey string
	err := h.DB.Get(&storedAPIKey, "SELECT api_key FROM device_info WHERE serial_number = $1", req.Info.Serial)
	if err != nil {
		log.Printf("DB error fetching device serial_number=%s: %v", req.Info.Serial, err)
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}

	// Validate API key
	if req.Info.APIKey != storedAPIKey {
		http.Error(w, "Unauthorized: invalid API key", http.StatusUnauthorized)
		return
	}

	// Insert telemetry data
	_, err = h.DB.Exec(`
		INSERT INTO telemetry (serial_number, vibration, x_accel, y_accel, z_accel)
		VALUES ($1, $2, $3, $4, $5)
	`, req.Info.Serial, req.Telemetry.Vibration, req.Telemetry.XAccel, req.Telemetry.YAccel, req.Telemetry.ZAccel)
	if err != nil {
		log.Printf("DB error inserting telemetry for serial_number=%s: %v", req.Info.Serial, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Update last_seen
	_, err = h.DB.Exec("UPDATE device_info SET last_seen = NOW() WHERE serial_number = $1", req.Info.Serial)
	if err != nil {
		log.Printf("DB error updating last_seen for serial_number=%s: %v", req.Info.Serial, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp := models.TelemetryResponse{
		Status:       true,
		SerialNumber: req.Info.Serial,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode telemetry response: %v", err)
	}
}