package mock

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathAuthenticate(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "authenticate",
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathAuthenticateRead,
			},
		},
	}
}

// pathAuthenticateRead is the endpoint Monzo will redirect to
func (b *backend) pathAuthenticateRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	// TODO: take the params from the redirect, and do the rest of the auth stuff
	// the stuff we need is in req.map.code and req.map.state, or data.Raw.code and data.Raw.state

	return &logical.Response{
		Data: map[string]interface{}{
			"ctx":  ctx,
			"req":  req,
			"data": data,
		},
	}, nil
}
