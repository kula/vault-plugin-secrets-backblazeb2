package b2

import (
    "context"
    "strings"

    "github.com/hashicorp/errwrap"
    "github.com/hashicorp/vault/logical"
    "github.com/hashicorp/vault/logical/framework"
)

// The stored plugin configuration

type Config struct {
    AccountId string `json:"account_id"`
    KeyId string `json:"key_id"`
    Key string `json:"key"`
    KeyName string `json:"key_name"`
    Configured bool `json:"is_configured"`
}

// A default, empty configuration

func DefaultConfig() *Config {
    return &Config{
	AccountId: "",
	KeyId: "",
	Key: "",
	KeyName: "vault-plugin-secrets-backblazeb2",
	Configured: false,
    }
}

// Update the configuration with new values

func (c *Config) Update(d *framework.FieldData) (bool, error) {
    if d == nil {
	return false, nil
    }

    changed := false

    keys := []string{"account_id", "key_id", "key", "key_name"}

    for _, key := range keys {
	if v, ok := d.GetOk(key); ok {
	    nv := strings.TrimSpace(v.(string))

	    switch key {
	    case "account_id":
		c.AccountId = nv
		c.Configured = true
		changed = true
	    case "key_id":
		c.KeyId = nv
		c.Configured = true
		changed = true
	    case "key":
		c.Key = nv
		c.Configured = true
		changed = true
	    case "key_name":
		c.KeyName = nv
		c.Configured = true
		changed = true
	    }
	}
    }

    // If you do not supply "key_id" it's assumed you're using your master
    // key, and key_id will be set the same as your account_id

    if c.KeyId == "" {
	c.KeyId = c.AccountId
    }

    return changed, nil
}


// Retrieve the configuration from the backend storage. In the event
// the plugin has not been configured, a default, empty configuration
// is returned

func (b *backend) GetConfig(ctx context.Context, s logical.Storage) (*Config, error) {
    c := DefaultConfig()

    entry, err := s.Get(ctx, "config");
    if err != nil {
	return nil, errwrap.Wrapf("failed to get configuration from backend: {{err}}", err)
    }

    if entry == nil || len(entry.Value) == 0 {
	return c, nil
    }

    if err := entry.DecodeJSON(&c); err != nil {
	return nil, errwrap.Wrapf("failed to decode configuration: {{err}}", err)
    }

    return c, nil
}
