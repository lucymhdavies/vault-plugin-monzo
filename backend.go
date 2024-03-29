package plugin

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// backend wraps the backend framework and adds a map for storing key value pairs
type backend struct {
	*framework.Backend

	// TODO: Monzo Client (with a mutex when reading/writing)
}

var _ logical.Factory = Factory

// Factory configures and returns Mock backends
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b, err := newBackend()
	if err != nil {
		return nil, err
	}

	if conf == nil {
		return nil, fmt.Errorf("configuration passed into backend is nil")
	}

	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}

	return b, nil

	// TODO: is this run every time Vault starts up, or only when the secret engine is first mounted?
	// if it's the latter, we'll need to find some way to instantiate the oauth2 client
	// when Vault restarts
}

func newBackend() (*backend, error) {
	b := backend{}

	b.Backend = &framework.Backend{
		Help:        strings.TrimSpace(mockHelp),
		BackendType: logical.TypeLogical,
		PathsSpecial: &logical.Paths{
			LocalStorage: []string{},
			SealWrapStorage: []string{
				"config",
			},

			// unauthenticated path for the redirect
			Unauthenticated: []string{
				"callback",
			},
		},
		Paths: framework.PathAppend(
			[]*framework.Path{
				pathConfig(&b),
				pathCallback(&b),
				pathAuthURL(&b),
			},
		),
	}

	return &b, nil
}

const mockHelp = `
The Monzo backend generates a Monzo API client and keeps its authentication token renewed
`
