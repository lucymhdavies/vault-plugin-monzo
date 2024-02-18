package plugin

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/oauth2"
)

const (
	tokenStoragePath = "token"
)

func pathToken(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "token",
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathTokenRead,
			},
		},
	}
}

// pathTokenRead returns the current oauth token, or an error if one does not exisst yet
func (b *backend) pathTokenRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	token, err := getToken(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	// TODO: if token is nil, return error?

	return &logical.Response{
		Data: map[string]interface{}{
			"access_token":  token.AccessToken,
			"expiry":        token.Expiry,
			"refresh_token": token.RefreshToken,
			"token_type":    token.TokenType,
		},
	}, nil
}

// getToken gets the current token from storage
func getToken(ctx context.Context, s logical.Storage) (*oauth2.Token, error) {
	entry, err := s.Get(ctx, tokenStoragePath)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	token := new(oauth2.Token)
	if err := entry.DecodeJSON(&token); err != nil {
		return nil, fmt.Errorf("error reading root configuration: %w", err)
	}

	// return the token, we are done
	return token, nil
}

// TODO: setToken function
