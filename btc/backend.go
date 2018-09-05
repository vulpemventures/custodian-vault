package btc

import (
	"context"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

type backend struct {
	*framework.Backend
}

// Factory ..
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

// Backend ..
func Backend(c *logical.BackendConfig) (*backend, error) {
	var b backend

	b.Backend = &framework.Backend{
		BackendType: logical.TypeLogical,
		Help:        BackendHelp,
		PathsSpecial: &logical.Paths{
			LocalStorage: []string{
				PathSecrets,
			},
		},
		Paths: []*framework.Path{
			pathWallet(&b),
			pathAddress(&b),
			pathCredentials(&b),
			pathTransaction(&b),
			pathMultiSigCredentials(&b),
			pathMultiSigWallet(&b),
			pathMultiSigAddress(&b),
		},
		Secrets: []*framework.Secret{
			secretCredentials(&b),
			multisigSecretCredentials(&b),
		},
	}

	return &b, nil
}
