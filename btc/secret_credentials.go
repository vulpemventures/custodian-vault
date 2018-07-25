package btc

import (
	"context"
	"time"
	"errors"

	"github.com/hashicorp/vault/helper/salt"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const SecretCredsType = "creds"

func secretCredentials(b *backend) *framework.Secret {
	return &framework.Secret{
		Type: SecretCredsType,
		Fields: map[string]*framework.FieldSchema{
			"token": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Access token for wallets",
			},
		},
		DefaultDuration: time.Duration(1 * time.Hour),
		Renew:  b.secretCredsRenew,
		Revoke: b.secretCredsRevoke,
	}
}

func (b *backend) secretCredsRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	resp := &logical.Response{Secret: req.Secret}

	resp.Secret.TTL = time.Duration(1 * time.Hour)
	resp.Secret.MaxTTL = time.Duration(2 * time.Hour)
	return resp, nil
}

func (b *backend) secretCredsRevoke(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	id, ok := req.Secret.InternalData["token"].(string)
	if !ok {
		return nil, errors.New("secret is missing internal data")
	}

	s, err := salt.NewSalt(ctx, req.Storage, nil)
	if err != nil {
		return nil, err
	}

	err = req.Storage.Delete(ctx, "creds/" + s.SaltID(id))
	if err != nil {
		return nil, err
	}

	return nil, nil
}