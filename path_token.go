package plugin

import (
	"context"
	"errors"
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

	if token == nil {
		return nil, errors.New("no token available")
	}

	// TODO: if we do have a token, but it has expired, clear it from storage

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

	// TODO: We should do a bunch of the same logic we currently have in renewToken
	// i.e. before naively returning the token, check if it needs a refresh first
	// however... renewToken currently calls getToken, so we'll need to move stuff around first

	// return the token, we are done
	return token, nil
}

// TODO: setToken function

// TODO: renewToken function will refresh the token if needed
func (b *backend) renewToken(ctx context.Context, req *logical.Request) error {
	b.Logger().Debug("renewToken called")

	config, err := getConfig(ctx, req.Storage)
	if err != nil {
		b.Logger().Error("getConfig err", err)
		return err
	}
	oauthConfig := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Scopes:       []string{}, // Monzo API has no documented scopes
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.AuthURL,
			TokenURL: config.TokenURL,
		},
		RedirectURL: config.RedirectBaseURL + "/v1/monzo/callback",
	}

	token, err := getToken(ctx, req.Storage)
	if err != nil {
		b.Logger().Error("getToken err", err)
		return err
	}

	if token == nil {
		b.Logger().Info("Token is currently nil")
		return nil
	}

	// Create a TokenSource that can auto-refresh
	ts := oauth2.ReuseTokenSource(token, oauthConfig.TokenSource(ctx, token))

	// For testing, force the token to refresh...
	// token.Expiry = time.Now().Add(-time.Minute)

	// Get the current (or refreshed) token
	current, err := ts.Token()
	if err != nil {
		b.Logger().Error("ts.Token() err", err)
		return err
	}

	b.Logger().Debug("Current token", "expiry", current.Expiry)

	// Persist if changed
	// TODO: call setToken
	if current.AccessToken != token.AccessToken {
		b.Logger().Debug("There was a new token", "token", current)

		// TODO: We probably want to modify the config somewhat, so it can renew before it actually expires

		entry, err := logical.StorageEntryJSON(tokenStoragePath, current)
		if err != nil {
			return err
		}
		if err := req.Storage.Put(ctx, entry); err != nil {
			return err
		}
		b.Logger().Info("New token persisted to storage")
	}

	return nil
}
