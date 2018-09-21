package btc

import (
	"context"
	"errors"
	"time"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/salt"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func secretCredentials(b *backend) *framework.Secret {
	return &framework.Secret{
		Type: SecretCredsType,
		Fields: map[string]*framework.FieldSchema{
			"token": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Access token for wallets",
			},
		},
		DefaultDuration: time.Duration(5 * time.Minute),
		Revoke:          b.secretCredsRevoke,
	}
}

func (b *backend) secretCredsRevoke(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	id, ok := req.Secret.InternalData["token"].(string)
	if !ok {
		return nil, errors.New(MissingInternalDataError)
	}

	s, err := salt.NewSalt(ctx, req.Storage, nil)
	if err != nil {
		return nil, err
	}

	err = req.Storage.Delete(ctx, PathCreds+s.SaltID(id))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func newToken(ctx context.Context, s logical.Storage, config *salt.Config) (string, string, error) {
	token, err := uuid.GenerateUUID()
	if err != nil {
		return "", "", err
	}

	newSalt, err := salt.NewSalt(ctx, s, config)
	if err != nil {
		return "", "", err
	}

	return token, newSalt.SaltID(token), nil
}

func (b *backend) GetToken(ctx context.Context, s logical.Storage, token string, walletType int) (*credential, error) {
	newSalt, err := salt.NewSalt(ctx, s, nil)
	if err != nil {
		return nil, err
	}

	leaseID := newSalt.SaltID(token)

	var path string
	switch walletType {
	case StandardType:
		path = PathCreds
	case MultiSigType:
		path = PathMultiSigCreds
	case SegWitType:
		path = PathSegWitCreds
	}

	path = path + leaseID

	entry, err := s.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var cred credential
	if err := entry.DecodeJSON(&cred); err != nil {
		return nil, err
	}

	return &cred, nil
}

func (b *backend) RevokeToken(ctx context.Context, store logical.Storage, token *credential, walletType int) error {
	var secretType string
	switch walletType {
	case StandardType:
		secretType = SecretCredsType
	case MultiSigType:
		secretType = MultiSigSecretCredsType
	case SegWitType:
		secretType = SegWitSecretCredsType
	}

	secret := b.Secret(secretType)
	request := &logical.Request{
		Storage:   store,
		Operation: logical.RevokeOperation,
		Secret: &logical.Secret{
			InternalData: map[string]interface{}{"token": token.Token},
			LeaseID:      token.LeaseID,
		},
	}
	_, err := secret.HandleRevoke(ctx, request)
	if err != nil {
		return err
	}

	return nil
}
