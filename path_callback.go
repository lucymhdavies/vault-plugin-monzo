package plugin

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/oauth2"
)

func pathCallback(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "callback",
		Fields: map[string]*framework.FieldSchema{
			"code": {
				Type:        framework.TypeString,
				Description: "The code returned from the Monzo redirect",
				Required:    true,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:      "Auth Code",
					Sensitive: false,
				},
			},
			"state": {
				Type:        framework.TypeString,
				Description: "The state returned from the Monzo redirect",
				Required:    true,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:      "Auth State",
					Sensitive: false,
				},
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathCallbackRead,
			},
		},
	}
}

// pathCallbackRead is the endpoint Monzo will redirect to
func (b *backend) pathCallbackRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// TODO: take the params from the redirect, and do the rest of the auth stuff
	// the stuff we need is in req.map.code and req.map.state, or data.Raw.code and data.Raw.state
	// Ideally actually use https://pkg.go.dev/github.com/hashicorp/vault/sdk/framework#FieldData

	config, err := getConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

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

	// TODO: Validate that we have a "code" in our payload
	var code string = data.Raw["code"].(string)

	// TODO: validate that the state is unchanged from what we generated in auth_url
	// this can probably be done by persisting the oauth client in its entirety
	//var state string = data.Raw["state"].(string)
	//oauthConfig.AuthCodeURL(state)

	b.Logger().Info("Attempting token exchange")

	tok, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		b.Logger().Error(err.Error())
		return nil, err
	}

	b.Logger().Debug("We got token", "token", tok)

	entry, err := logical.StorageEntryJSON(tokenStoragePath, tok)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}
	b.Logger().Info("Token persisted to storage")

	return &logical.Response{
		Data: map[string]interface{}{
			"message": "Authenticated. Log in to the Monzo app to grant permissions",
		},
	}, nil
}
