package config

import (
	"os"
	"testing"
)

func TestLoadVaultConfig(t *testing.T) {
	os.Setenv("VAULT_ADDR", "http://127.0.0.1:8200")
	os.Setenv("VAULT_TOKEN", "hvs.testtoken")

	cfg := LoadVaultConfig()

	if cfg.VaultAddr != "http://127.0.0.1:8200" {
		t.Errorf("expected VaultAddr to be http://127.0.0.1:8200, got %s", cfg.VaultAddr)
	}
	if cfg.VaultToken != "hvs.testtoken" {
		t.Errorf("expected VaultToken to be hvs.testtoken, got %s", cfg.VaultToken)
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

func TestWithDBCredentials(t *testing.T) {
	cfg := Config{}

	creds := map[string]string{
		"host":     "localhost",
		"port":     "5432",
		"user":     "bankii",
		"password": "NYAJOWI",
		"dbname":   "techtarget_project",
	}

	cfg = cfg.WithDBCredentials(creds)

	if cfg.DBHost != "localhost" {
		t.Errorf("expected DBHost localhost, got %s", cfg.DBHost)
	}
	if cfg.DBPass != "NYAJOWI" {
		t.Errorf("expected DBPass NYAJOWI, got %s", cfg.DBPass)
	}
}