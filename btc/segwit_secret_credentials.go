package btc

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/vault/helper/salt"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func segwitSecretCredentials(b *backend) *framework.Secret {
	return &framework.Secret{
		Type: SegWitSecretCredsType,
		Fields: map[string]*framework.FieldSchema{
			"token": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Access token for wallets",
			},
		},
		DefaultDuration: time.Duration(5 * time.Minute),
		Revoke:          b.segwitSecretCredsRevoke,
	}
}

func (b *backend) segwitSecretCredsRevoke(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	id, ok := req.Secret.InternalData["token"].(string)
	if !ok {
		return nil, errors.New(MissingInternalDataError)
	}

	s, err := salt.NewSalt(ctx, req.Storage, nil)
	if err != nil {
		return nil, err
	}

	err = req.Storage.Delete(ctx, PathSegWitCreds+s.SaltID(id))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
