package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type FirmwareUpdateResponse struct {
	Status              int    `json:"status"`
	CurrentFirmware     string `json:"current_firmware"`
	NextFirmwareVersion string `json:"next_firmware_version"`
}

var db *sqlx.DB

func main() {
	var err error

	connStr := "host=localhost port=5432 user=bankii password=NYAJOWI dbname=techtarget_project sslmode=disable"

	// OPEN CONNECTION
	db, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalln("Failed to connect to database:", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatalln("Database is unreachable:", err)
	}

	fmt.Println("Successfully connected to techtarget_project!")

	// REGISTER  HANDLER
	http.HandleFunc("/check-firmware", CheckFirmwareHandler)

	// START SERVER
	fmt.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

//
func CheckFirmwareHandler(w http.ResponseWriter, r *http.Request) {
	receivedVersion := r.URL.Query().Get("v") // Example: /check-firmware?v=v2.34

	var dbVersion string
	// Fetching from your device_info table
	err := db.Get(&dbVersion, "SELECT firmware FROM device_info LIMIT 1")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if receivedVersion != dbVersion {
		resp := FirmwareUpdateResponse{
			Status:              200,
			CurrentFirmware:     receivedVersion,
			NextFirmwareVersion: dbVersion,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	} else {
		// If they match, you can send a simple 'OK' message
		w.Write([]byte("Firmware is up to date"))
	}
}