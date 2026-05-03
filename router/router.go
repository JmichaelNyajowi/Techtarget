package router

import (
	"net/http"

	"techtarget_project/handlers"
	"techtarget_project/vault"

	"github.com/jmoiron/sqlx"
)

func Setup(db *sqlx.DB, vaultClient *vault.Client, universalKey string) *http.ServeMux {
	mux := http.NewServeMux()

	firmwareHandler := &handlers.FirmwareHandler{DB: db}
	mux.HandleFunc("/check-firmware", firmwareHandler.CheckFirmware)

	telemetryHandler := &handlers.TelemetryHandler{DB: db}
	mux.HandleFunc("/telemetry", telemetryHandler.Submit)

	provisionHandler := &handlers.ProvisionHandler{
		DB:           db,
		VaultClient:  vaultClient,
		UniversalKey: universalKey,
	}
	mux.HandleFunc("/provision", provisionHandler.Provision)

	return mux
}