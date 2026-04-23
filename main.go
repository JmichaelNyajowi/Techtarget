package main

import (
    "fmt"
    "log"
    "net/http"
    "time"

    "techtarget_project/config"
    "techtarget_project/db"
    "techtarget_project/router"
)

func main() {
    cfg := config.Load()

    database := db.Connect(cfg.DSN())
    defer database.Close()

    mux := router.Setup(database)

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