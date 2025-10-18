package plugin

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

func pathAuthURL(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "auth-url",
		Fields: map[string]*framework.FieldSchema{
			"state": {
				Type:        framework.TypeString,
				Description: "A random state to use in the auth",
				Required:    false,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:      "Auth State",
					Sensitive: false,
				},
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.pathAuthURLWrite,
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathAuthURLWrite,
			},
		},
	}
}

// pathAuthURLWrite requests a new Auth URL from Vault
func (b *backend) pathAuthURLWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := getConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	// TODO: check if we've got a valid config yet, and fail if we do not

	oauthConfig := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Scopes:       []string{}, // Monzo API has no documented scopes
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.AuthURL,
			TokenURL: config.TokenURL,
		},

		// TODO: figure this out automatically from req? May not be possible
		RedirectURL: config.RedirectBaseURL + "/v1/monzo/callback",
	}

	url := oauthConfig.AuthCodeURL(uuid.New().String())
	// TODO: We probably want to persist this specific oauth2 client...

	return &logical.Response{
		Data: map[string]interface{}{
			"url": url,
		},
	}, nil
}
