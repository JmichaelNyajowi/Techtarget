package db

import (
    "log"
    "time"

    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

func Connect(dsn string) *sqlx.DB {
    database, err := sqlx.Connect("postgres", dsn)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    database.SetMaxOpenConns(10)
    database.SetMaxIdleConns(5)
    database.SetConnMaxLifetime(5 * time.Minute)

    return database
}