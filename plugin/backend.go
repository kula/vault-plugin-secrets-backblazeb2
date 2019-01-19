package b2

import (
    "context"

    "github.com/hashicorp/vault/logical"
    "github.com/hashicorp/vault/logical/framework"
)

type backend struct {
    *framework.Backend
}

// Factory returns a configured instance of the B2 backend
func Factory(ctx context.Context, c *logical.BackendConfig) (logical.Backend, error) {
    b := Backend()
    if err := b.Setup(ctx, c); err != nil {
	return nil, err
    }

    b.Logger().Info("Plugin successfully initialized")
    return b, nil
}

// Backend returns a configured B2 backend
func Backend() *backend {
    var b backend

    b.Backend = &framework.Backend{
	BackendType: logical.TypeLogical,
	Help: "The B2 secrets backend provisions API keys for the Backblaze B2 service",

	Paths: []*framework.Path{
	    // path_config.go
	    // ^config
	    b.pathConfigCRUD(),
	},
    }

    return &b
}
