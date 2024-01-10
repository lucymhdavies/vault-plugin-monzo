package mock

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	configStoragePath = "config"
)

// monzoConfig includes the minimum configuration
// required to instantiate a new Monzo client.
type monzoConfig struct {
	ClientID     string `json:"username"`
	ClientSecret string `json:"password"`
	AuthURL      string `json:"auth_url"`
	TokenURL     string `json:"token_url"`

	// TODO: what if we just had an oauth2.Config ?
}

// pathConfig extends the Vault API with a `/config`
// endpoint for the backend. You can choose whether
// or not certain attributes should be displayed,
// required, and named. For example, password
// is marked as sensitive and will not be output
// when you read the configuration.
func pathConfig(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config",
		Fields: map[string]*framework.FieldSchema{
			"client_id": {
				Type:        framework.TypeString,
				Description: "OAuth Client ID",
				Required:    true,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:      "Client ID",
					Sensitive: false,
				},
			},
			"client_secret": {
				Type:        framework.TypeString,
				Description: "OAuth Client Secret",
				Required:    true,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:      "Client Secret",
					Sensitive: true,
				},
			},
			"auth_url": {
				Type:        framework.TypeString,
				Description: "The Monzo Auth URL",
				Required:    false,
				Default:     "https://auth.monzo.com/",
				DisplayAttrs: &framework.DisplayAttributes{
					Name:      "Auth URL",
					Sensitive: false,
				},
			},
			"token_url": {
				Type:        framework.TypeString,
				Description: "The Monzo Token URL",
				Required:    false,
				Default:     "https://api.monzo.com/oauth2/token",
				DisplayAttrs: &framework.DisplayAttributes{
					Name:      "Token URL",
					Sensitive: false,
				},
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathConfigRead,
			},
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.pathConfigWrite,
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathConfigWrite,
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.pathConfigDelete,
			},
		},
		ExistenceCheck:  b.pathConfigExistenceCheck,
		HelpSynopsis:    pathConfigHelpSynopsis,
		HelpDescription: pathConfigHelpDescription,
	}
}

// pathConfigExistenceCheck verifies if the configuration exists.
func (b *backend) pathConfigExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	out, err := req.Storage.Get(ctx, req.Path)
	if err != nil {
		return false, fmt.Errorf("existence check failed: %w", err)
	}

	return out != nil, nil
}

// pathConfigRead reads the configuration and outputs non-sensitive information.
func (b *backend) pathConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := getConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"client_id": config.ClientID,
			"auth_url":  config.AuthURL,
			"token_url": config.TokenURL,
		},
	}, nil
}

// pathConfigWrite updates the configuration for the backend
func (b *backend) pathConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := getConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	createOperation := (req.Operation == logical.CreateOperation)

	if config == nil {
		if !createOperation {
			return nil, errors.New("config not found during update operation")
		}
		config = new(monzoConfig)
	}

	// set defaults (on either create or update)
	config.AuthURL = "https://auth.monzo.com/"
	config.TokenURL = "https://api.monzo.com/oauth2/token"

	if clientID, ok := data.GetOk("client_id"); ok {
		config.ClientID = clientID.(string)
	} else if !ok && createOperation {
		return nil, fmt.Errorf("missing client_id in configuration")
	}

	if clientSecret, ok := data.GetOk("client_secret"); ok {
		config.ClientSecret = clientSecret.(string)
	} else if !ok && createOperation {
		return nil, fmt.Errorf("missing client_secret in configuration")
	}

	if authURL, ok := data.GetOk("auth_url"); ok {
		config.AuthURL = authURL.(string)
	}

	if tokenURL, ok := data.GetOk("token_url"); ok {
		config.TokenURL = tokenURL.(string)
	}

	entry, err := logical.StorageEntryJSON(configStoragePath, config)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	// TODO: Create a new Monzo client, based on
	// https://github.com/lucymhdavies/monzo-token-renewer/blob/main/main.go#L74-L94

	b.Logger().Info("New Monzo Client Created")

	return nil, nil
}

// pathConfigDelete removes the configuration for the backend
func (b *backend) pathConfigDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete(ctx, configStoragePath)

	//if err == nil {
	//b.reset()
	//}

	return nil, err
}

func getConfig(ctx context.Context, s logical.Storage) (*monzoConfig, error) {
	entry, err := s.Get(ctx, configStoragePath)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	config := new(monzoConfig)
	if err := entry.DecodeJSON(&config); err != nil {
		return nil, fmt.Errorf("error reading root configuration: %w", err)
	}

	// return the config, we are done
	return config, nil
}

// pathConfigHelpSynopsis summarizes the help text for the configuration
const pathConfigHelpSynopsis = `Configure the HashiCups backend.`

// pathConfigHelpDescription describes the help text for the configuration
const pathConfigHelpDescription = `
The HashiCups secret backend requires credentials for managing
JWTs issued to users working with the products API.

You must sign up with a username and password and
specify the HashiCups address for the products API
before using this secrets backend.
`
