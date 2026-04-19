package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type FirmwareUpdateResponse struct {
	Status              string `json:"status"`
	CurrentFirmware     string `json:"current_firmware"`
	NextFirmwareVersion string `json:"next_firmware_version"`
}

var db *sqlx.DB

func main() {
	var err error

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
	)

	db, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	mux := http.NewServeMux()
	mux.HandleFunc("/check-firmware", checkFirmwareHandler)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(srv.ListenAndServe())
}

func checkFirmwareHandler(w http.ResponseWriter, r *http.Request) {
	// Enforce GET only
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Validate required query parameters
	receivedVersion := r.URL.Query().Get("v")
	if receivedVersion == "" {
		http.Error(w, "Missing required query parameter: v", http.StatusBadRequest)
		return
	}

	deviceID := r.URL.Query().Get("serial_number")
	if deviceID == "" {
		http.Error(w, "Missing required query parameter: serial_number", http.StatusBadRequest)
		return
	}

	// Query firmware for the specific device
	var dbVersion string
	err := db.Get(&dbVersion, "SELECT firmware FROM device_info WHERE serial_number = $1", deviceID)
	if err != nil {
		log.Printf("DB error for serial_number=%s: %v", deviceID, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if receivedVersion != dbVersion {
		// Update available
		w.WriteHeader(http.StatusAccepted) // 202
		resp := FirmwareUpdateResponse{
			Status:              "update_available",
			CurrentFirmware:     receivedVersion,
			NextFirmwareVersion: dbVersion,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Failed to encode response: %v", err)
		}
	} else {
		// Already up to date
		w.WriteHeader(http.StatusOK) // 200
		resp := map[string]string{
			"status":           "up_to_date",
			"current_firmware": receivedVersion,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Failed to encode response: %v", err)
		}
	}
}
