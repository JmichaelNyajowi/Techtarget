package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"techtarget_project/config"
	"techtarget_project/db"
	"techtarget_project/router"
	"techtarget_project/vault"
)

func main() {
	// Load Vault connection details from environment
	cfg := config.LoadVaultConfig()

	// Connect to Vault
	vaultClient, err := vault.New(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		log.Fatalf("Failed to connect to Vault: %v", err)
	}

	// Fetch DB credentials from Vault
	creds, err := vaultClient.GetDBCredentials(cfg.VaultDBPath)
	if err != nil {
		log.Fatalf("Failed to fetch DB credentials from Vault: %v", err)
	}
	cfg = cfg.WithDBCredentials(creds)

	// Fetch universal provisioning key from Vault
	universalKey, err := vaultClient.GetUniversalKey(cfg.VaultUniversalPath)
	if err != nil {
		log.Fatalf("Failed to fetch universal key from Vault: %v", err)
	}

	// Connect to database
	database := db.Connect(cfg.DSN())
	defer database.Close()

	// Set up router
	mux := router.Setup(database, vaultClient, universalKey)

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