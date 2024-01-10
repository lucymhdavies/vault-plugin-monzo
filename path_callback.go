package plugin

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
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

	return &logical.Response{
		Data: map[string]interface{}{
			"data": data,
		},
	}, nil
}
