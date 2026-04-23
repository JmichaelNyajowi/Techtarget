package router

import (
    "net/http"

    "techtarget_project/handlers"

    "github.com/jmoiron/sqlx"
)

func Setup(db *sqlx.DB) *http.ServeMux {
    mux := http.NewServeMux()

    firmwareHandler := &handlers.FirmwareHandler{DB: db}
    mux.HandleFunc("/check-firmware", firmwareHandler.CheckFirmware)

    telemetryHandler := &handlers.TelemetryHandler{DB: db}
	mux.HandleFunc("/telemetry", telemetryHandler.Submit)

    return mux
}