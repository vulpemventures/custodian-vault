package btc

import (
	"context"
	"errors"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathSegWitAddress(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: PathSegWitAddress + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Wallet name",
			},
			"token": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Auth token",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathSegWitAddressWrite,
		},

		HelpSynopsis:    PathSegWitAddressHelpSyn,
		HelpDescription: PathSegWitAddressHelpDesc,
	}
}

func (b *backend) pathSegWitAddressWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	walletName := d.Get("name").(string)
	if walletName == "" {
		return nil, errors.New(MissingWalletNameError)
	}
	t := d.Get("token").(string)
	if t == "" {
		return nil, errors.New(MissingTokenError)
	}
	// add prefix for segwit wallet
	walletName = SegWitPrefix + walletName

	// check if auth token is valid
	token, err := b.GetToken(ctx, req.Storage, t, SegWitType)
	if err != nil {
		return nil, err
	}
	if token == nil || walletName != token.WalletName {
		return nil, errors.New(InvalidTokenError)
	}

	w, err := b.GetSegWitWallet(ctx, req.Storage, walletName)
	if err != nil {
		return nil, err
	}

	childnum, err := b.GetLastUsedAddressIndex(ctx, req.Storage, walletName, true)
	if err != nil {
		return nil, err
	}

	// increment childnum to derive next address
	childnum = childnum + 1
	a, err := deriveSegWitAddress(w, childnum)
	if err != nil {
		return nil, err
	}

	// override the storage with new generated address
	entry, err := logical.StorageEntryJSON(PathSegWitAddress+walletName, a)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	// revoke auth token
	err = b.RevokeToken(ctx, req.Storage, token, SegWitType)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"address": a.LastAddress,
		},
	}, nil
}
