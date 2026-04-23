package config

import (
    "fmt"
    "os"
)

type Config struct {
    DBHost string
    DBPort string
    DBUser string
    DBPass string
    DBName string
}

func Load() Config {
    return Config{
        DBHost: os.Getenv("DB_HOST"),
        DBPort: os.Getenv("DB_PORT"),
        DBUser: os.Getenv("DB_USER"),
        DBPass: os.Getenv("DB_PASS"),
        DBName: os.Getenv("DB_NAME"),
    }
}

func (c Config) DSN() string {
    return fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        c.DBHost, c.DBPort, c.DBUser, c.DBPass, c.DBName,
    )
}