package btc

import (
	"context"
	"errors"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathSegWitWallet(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: PathSegWitWallet + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"network": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "btc network type: mainnet | testnet",
			},
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "wallet name",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathSegWitWalletRead,
			logical.UpdateOperation: b.pathSegWitWalletWrite,
		},

		HelpSynopsis:    PathSegWitWalletsHelpSyn,
		HelpDescription: PathSegWitWalletsHelpDesc,
	}
}

func (b *backend) pathSegWitWalletWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	network := d.Get("network").(string)
	if network == "" {
		return nil, errors.New(MissingNetworkError)
	}
	if network != "testnet" && network != "mainnet" {
		return nil, errors.New(InvalidNetworkError)
	}

	walletName := d.Get("name").(string)
	if walletName == "" {
		return nil, errors.New(MissingWalletNameError)
	}

	// return error if a wallet with same name has already been created
	w, err := b.GetSegWitWallet(ctx, req.Storage, walletName)
	if err != nil {
		return nil, err
	}
	if w != nil {
		return nil, errors.New(SegWitWalletAlreadyExistsError)
	}

	// create a new wallet
	wallet, err := createSegWitWallet(network)
	if err != nil {
		return nil, err
	}

	// create storage entry
	entry, err := logical.StorageEntryJSON(PathSegWitWallet+walletName, wallet)
	if err != nil {
		return nil, err
	}

	// save in local storage
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"mnemonic": wallet.Mnemonic,
		},
	}, nil
}

func (b *backend) pathSegWitWalletRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	walletName := d.Get("name").(string)

	// get wallet from storage
	w, err := b.GetSegWitWallet(ctx, req.Storage, walletName)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return nil, nil
	}

	seed := seedFromMnemonic(w.Mnemonic)

	// get master key from seed
	key, err := getMasterKey(seed, w.Network)
	if err != nil {
		return nil, err
	}

	// first derive private key at path m/84'/0'/0'/0 (mainnet)
	xprv, err := derivePrivKey(key, w.DerivationPath)
	if err != nil {
		return nil, err
	}

	// then derive public key and return in xpub format
	xpub, err := derivePubKey(xprv)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"network": w.Network,
			"xpub":    xpub.String(),
		},
	}, nil
}

// Retrieves a wallet in storage given the wallet name
func (b *backend) GetSegWitWallet(ctx context.Context, store logical.Storage, walletName string) (*wallet, error) {
	entry, err := store.Get(ctx, PathSegWitWallet+walletName)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var w wallet
	if err := entry.DecodeJSON(&w); err != nil {
		return nil, err
	}

	return &w, nil
}
