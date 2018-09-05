package btc

import (
	"context"
	"errors"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

type multiSigWallet struct {
	Network        string
	Mnemonic       string
	DerivationPath []uint32
	M              int
	N              int
	RedeemScript   string
	PublicKeys     []string
}

func pathMultiSigWallet(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: PathMultiSigWallet + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"network": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Btc network type: mainnet | testnet",
			},
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Wallet name",
			},
			"pubkeys": &framework.FieldSchema{
				Type:        framework.TypeCommaStringSlice,
				Description: "List of public keys for multisig wallet",
			},
			"m": &framework.FieldSchema{
				Type:        framework.TypeInt,
				Description: "Threshold signature",
			},
			"n": &framework.FieldSchema{
				Type:        framework.TypeInt,
				Description: "Total number of signatures",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathMultiSigWalletRead,
			logical.UpdateOperation: b.pathMultiSigWalletWrite,
		},

		HelpSynopsis:    PathMultiSigWalletsHelpSyn,
		HelpDescription: PathMultiSigWalletsHelpDesc,
	}
}

func (b *backend) pathMultiSigWalletWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	network := d.Get("network").(string)
	if network == "" {
		return nil, errors.New(MissingNetworkError)
	}

	walletName := d.Get("name").(string)
	if walletName == "" {
		return nil, errors.New(MissingWalletNameError)
	}

	pubkeys := d.Get("pubkeys").([]string)
	if len(pubkeys) == 0 {
		return nil, errors.New(MissingPubKeysError)
	}

	m := d.Get("m").(int)
	if m <= 0 {
		return nil, errors.New(InvalidMError)
	}

	n := d.Get("n").(int)
	if n <= 0 {
		return nil, errors.New(InvalidNError)
	}

	// check valid params:
	// # of public keys should be equal to n - 1
	// m should be minor or equal to n
	// TODO: check valid public keys
	if l := len(pubkeys); l != (n - 1) {
		return nil, errors.New("Invalid number of public keys: provided " + string(l) + " expected " + string(n-1))
	}
	if m > n {
		return nil, errors.New(MBiggerThanNError)
	}

	// return error if a wallet with same name has already been created
	walletName = MultiSigPrefix + walletName
	w, err := b.GetMultiSigWallet(ctx, req.Storage, walletName)
	if err != nil {
		return nil, err
	}
	if w != nil {
		return nil, errors.New("MultiSig wallet with name '" + walletName + "' already exists")
	}

	// create multisig wallet with params
	wallet, err := createMultiSigWallet(network, pubkeys, m, n)
	if err != nil {
		return nil, err
	}

	// create storage entry
	entry, err := logical.StorageEntryJSON("wallet/"+walletName, wallet)
	if err != nil {
		return nil, err
	}

	// save in local storage
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathMultiSigWalletRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	walletName := d.Get("name").(string)
	walletName = MultiSigPrefix + walletName

	// get wallet from storage
	w, err := b.GetMultiSigWallet(ctx, req.Storage, walletName)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"m":            w.M,
			"n":            w.N,
			"pubkeys":      w.PublicKeys,
			"redeemScript": w.RedeemScript,
		},
	}, nil
}

func (b *backend) GetMultiSigWallet(ctx context.Context, store logical.Storage, walletName string) (*multiSigWallet, error) {
	entry, err := store.Get(ctx, "wallet/"+walletName)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var w multiSigWallet
	if err := entry.DecodeJSON(&w); err != nil {
		return nil, err
	}

	return &w, nil
}
