package btc

import (
	"context"
	"errors"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathMultiSigAddress(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: PathMultiSigAddress + framework.GenericNameRegex("name"),
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
			logical.UpdateOperation: b.pathMultiSigAddressWrite,
		},

		HelpSynopsis:    PathMultiSigAddressHelpSyn,
		HelpDescription: PathMultiSigAddressHelpDesc,
	}
}

func (b *backend) pathMultiSigAddressWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	walletName := d.Get("name").(string)
	if walletName == "" {
		return nil, errors.New(MissingWalletNameError)
	}
	t := d.Get("token").(string)
	if t == "" {
		return nil, errors.New(MissingTokenError)
	}

	// check if auth token is valid
	token, err := b.GetToken(ctx, req.Storage, t, MultiSigType)
	if err != nil {
		return nil, err
	}
	if token == nil || walletName != token.WalletName {
		return nil, errors.New(InvalidTokenError)
	}

	w, err := b.GetMultiSigWallet(ctx, req.Storage, walletName)
	if err != nil {
		return nil, err
	}

	// for multisig, address is always the same and it's built from redeem script
	address, err := getMultiSigAddress(w.RedeemScript, w.Network)
	if err != nil {
		return nil, err
	}

	// revoke auth token
	err = b.RevokeToken(ctx, req.Storage, token, MultiSigType)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"address": address,
		},
	}, nil
}
