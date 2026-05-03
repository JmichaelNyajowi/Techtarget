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

	VaultAddr        string
	VaultToken       string
	VaultDBPath      string
	VaultUniversalPath string
}

func LoadVaultConfig() Config {
	return Config{
		VaultAddr:          os.Getenv("VAULT_ADDR"),
		VaultToken:         os.Getenv("VAULT_TOKEN"),
		VaultDBPath:        "secret/data/techtarget_project/db",
		VaultUniversalPath: "secret/data/techtarget_project/universal",
	}
}

func (c Config) WithDBCredentials(creds map[string]string) Config {
	c.DBHost = creds["host"]
	c.DBPort = creds["port"]
	c.DBUser = creds["user"]
	c.DBPass = creds["password"]
	c.DBName = creds["dbname"]
	return c
}

func (c Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPass, c.DBName,
	)
}