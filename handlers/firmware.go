package handlers

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/jmoiron/sqlx"

    "techtarget_project/models"
)

type FirmwareHandler struct {
    DB *sqlx.DB
}

func (h *FirmwareHandler) CheckFirmware(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    receivedVersion := r.URL.Query().Get("v")
    if receivedVersion == "" {
        http.Error(w, "Missing required query parameter: v", http.StatusBadRequest)
        return
    }

    serialNumber := r.URL.Query().Get("serial_number")
    if serialNumber == "" {
        http.Error(w, "Missing required query parameter: serial_number", http.StatusBadRequest)
        return
    }

    var dbVersion string
    err := h.DB.Get(&dbVersion, "SELECT firmware FROM device_info WHERE serial_number = $1", serialNumber)
    if err != nil {
        log.Printf("DB error for serial_number=%s: %v", serialNumber, err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")

    if receivedVersion != dbVersion {
        w.WriteHeader(http.StatusAccepted)
        resp := models.FirmwareUpdateResponse{
            Status:              false,
            CurrentFirmware:     receivedVersion,
            NextFirmwareVersion: dbVersion,
        }
        if err := json.NewEncoder(w).Encode(resp); err != nil {
            log.Printf("Failed to encode response: %v", err)
        }
    } else {
        w.WriteHeader(http.StatusOK)
        resp := models.FirmwareUpToDateResponse{
            Status:          true,
            CurrentFirmware: receivedVersion,
        }
        if err := json.NewEncoder(w).Encode(resp); err != nil {
            log.Printf("Failed to encode response: %v", err)
        }
    }
}