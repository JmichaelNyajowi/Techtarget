package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Set environment variables for the test
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "bankii")
	os.Setenv("DB_PASS", "NYAJOWI")
	os.Setenv("DB_NAME", "techtarget_project")

	cfg := Load()

	if cfg.DBHost != "localhost" {
		t.Errorf("expected DBHost to be localhost, got %s", cfg.DBHost)
	}
	if cfg.DBPort != "5432" {
		t.Errorf("expected DBPort to be 5432, got %s", cfg.DBPort)
	}
	if cfg.DBUser != "bankii" {
		t.Errorf("expected DBUser to be bankii, got %s", cfg.DBUser)
	}
}

func TestDSN(t *testing.T) {
	cfg := Config{
		DBHost: "localhost",
		DBPort: "5432",
		DBUser: "bankii",
		DBPass: "NYAJOWI",
		DBName: "techtarget_project",
	}

	expected := "host=localhost port=5432 user=bankii password=NYAJOWI dbname=techtarget_project sslmode=disable"
	got := cfg.DSN()

	if got != expected {
		t.Errorf("expected DSN:\n%s\ngot:\n%s", expected, got)
	}
}