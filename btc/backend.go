package btc

import (
	"context"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

type backend struct {
	*framework.Backend
}

func Factory(ctx context.Context, c *logical.BackendConfig) (logical.Backend, error) {
	b, err := Backend(c)
	if err != nil {
		return nil, err
	}

	if err := b.Setup(ctx, c); err != nil {
		return nil, err
	}

	return b, nil
}

func Backend(c *logical.BackendConfig) (*backend, error) {
	var b backend

	b.Backend = &framework.Backend{
		BackendType: logical.TypeLogical,
		Help: backendHelp,
		PathsSpecial: &logical.Paths{
			LocalStorage: []string{
				"secrets/",
			},
		},
		Paths: []*framework.Path{
			pathWallet(&b),
			pathAddress(&b),
			pathCredentials(&b),
		},
		Secrets: []*framework.Secret{
			secretCredentials(&b),
		},
	}

	return &b, nil
}

const backendHelp = `
The bitcoin custodian plugin lets you to store your private keys securely.
With this plugin you can create wallet private keys, generate receiving addresses and sign transactions.
`