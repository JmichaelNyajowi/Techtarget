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

	// Decode JSON body
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
	if req.Info.Model == "" {
		http.Error(w, "info.model is required", http.StatusBadRequest)
		return
	}
	if req.Info.Firmware == "" {
		http.Error(w, "info.firmware is required", http.StatusBadRequest)
		return
	}

	// Upsert device_info — insert if new device, update if existing
	_, err := h.DB.Exec(`
    INSERT INTO device_info (serial_number, model, firmware, last_seen)
    VALUES ($1, $2, $3, NOW())
    ON CONFLICT (serial_number)
    DO UPDATE SET 
        model = EXCLUDED.model, 
        firmware = EXCLUDED.firmware,
        last_seen = NOW()
`, req.Info.Serial, req.Info.Model, req.Info.Firmware)
	if err != nil {
		log.Printf("DB error upserting device_info for serial_number=%s: %v", req.Info.Serial, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp := models.TelemetryResponse{
		Status:       true,
		SerialNumber: req.Info.Serial,
		Message:      "Telemetry recorded and device info updated",
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}