package vault

import (
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
		// KV v2 wraps data inside a "data" key
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