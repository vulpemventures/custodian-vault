package btc

import (
	"context"
	"errors"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathMultiSigCredentials(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: PathMultiSigCreds + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Wallet name",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathMultiSigCredsRead,
		},

		HelpSynopsis:    PathMultiSigCredsHelpSyn,
		HelpDescription: PathMultiSigCredsHelpDesc,
	}
}

func (b *backend) pathMultiSigCredsRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	walletName := d.Get("name").(string)
	if walletName == "" {
		return nil, errors.New(MissingWalletNameError)
	}

	w, err := b.GetMultiSigWallet(ctx, req.Storage, walletName)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return nil, errors.New(MultiSigWalletNotFoundError)
	}

	token, leaseID, err := newToken(ctx, req.Storage, nil)
	if err != nil {
		return nil, err
	}

	cred := &credential{
		WalletName: walletName,
		LeaseID:    leaseID,
		Token:      token,
	}

	entry, err := logical.StorageEntryJSON(PathMultiSigCreds+leaseID, cred)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	resp := b.Secret(SecretCredsType).Response(
		map[string]interface{}{"token": token},
		map[string]interface{}{"token": token},
	)

	return resp, nil
}
