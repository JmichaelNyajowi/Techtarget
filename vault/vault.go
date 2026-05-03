package vault

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	vault "github.com/hashicorp/vault/api"
)

type Client struct {
	logical *vault.Logical
}

func New(address, token string) (*Client, error) {
	config := vault.DefaultConfig()
	config.Address = address

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	client.SetToken(token)

	return &Client{logical: client.Logical()}, nil
}

func (c *Client) GetDBCredentials(path string) (map[string]string, error) {
	secret, err := c.logical.Read(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret at %s: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no secret found at path: %s", path)
	}

	creds := make(map[string]string)
	for key, value := range secret.Data {
		if key == "data" {
			if nested, ok := value.(map[string]interface{}); ok {
				for k, v := range nested {
					creds[k] = fmt.Sprintf("%v", v)
				}
			}
		} else {
			creds[key] = fmt.Sprintf("%v", value)
		}
	}

	return creds, nil
}

// GetUniversalKey fetches the universal provisioning key from Vault
func (c *Client) GetUniversalKey(path string) (string, error) {
	secret, err := c.logical.Read(path)
	if err != nil {
		return "", fmt.Errorf("failed to read universal key at %s: %w", path, err)
	}
	if secret == nil {
		return "", fmt.Errorf("no secret found at path: %s", path)
	}

	// handle KV v2 nested data
	if data, ok := secret.Data["data"]; ok {
		if nested, ok := data.(map[string]interface{}); ok {
			if key, ok := nested["universal_key"]; ok {
				return fmt.Sprintf("%v", key), nil
			}
		}
	}

	return "", fmt.Errorf("universal_key not found in vault path: %s", path)
}

// GenerateAPIKey generates a cryptographically secure unique API key
func GenerateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate api key: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}