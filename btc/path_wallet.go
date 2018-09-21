package btc

import (
	"context"
	"errors"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

type wallet struct {
	Network        string
	Mnemonic       string
	DerivationPath []uint32
	Segwit         bool
}

func pathWallet(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: PathWallet + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"network": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "btc network type: mainnet | testnet | regtest",
			},
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "wallet name",
			},
			"segwit": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Description: "bip49 segwit backward compatible wallet",
				Default:     false,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathWalletRead,
			logical.UpdateOperation: b.pathWalletWrite,
		},

		HelpSynopsis:    PathWalletsHelpSyn,
		HelpDescription: PathWalletsHelpDesc,
	}
}

func (b *backend) pathWalletWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	network := d.Get("network").(string)
	if network == "" {
		return nil, errors.New(MissingNetworkError)
	}
	if network != MainNet && network != TestNet && network != RegTest {
		return nil, errors.New(InvalidNetworkError)
	}

	walletName := d.Get("name").(string)
	if walletName == "" {
		return nil, errors.New(MissingWalletNameError)
	}

	segwit := d.Get("segwit").(bool)

	// return error if a wallet with same name has already been created
	w, err := b.GetWallet(ctx, req.Storage, walletName)
	if err != nil {
		return nil, err
	}
	if w != nil {
		return nil, errors.New(WalletAlreadyExistsError)
	}

	// create a new wallet
	wallet, err := createWallet(network, segwit)
	if err != nil {
		return nil, err
	}

	// create storage entry
	entry, err := logical.StorageEntryJSON(PathWallet+walletName, wallet)
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

func (b *backend) pathWalletRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	walletName := d.Get("name").(string)

	// get wallet from storage
	w, err := b.GetWallet(ctx, req.Storage, walletName)
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

	// first derive private key at path m/44'/0'/0'/0
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
func (b *backend) GetWallet(ctx context.Context, store logical.Storage, walletName string) (*wallet, error) {
	entry, err := store.Get(ctx, PathWallet+walletName)
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

func getWalletByType(ctx context.Context, b *backend, store logical.Storage, walletName string, walletType int) (*wallet, error) {
	switch walletType {
	case StandardType:
		return b.GetWallet(ctx, store, walletName)
	case SegWitType:
		return b.GetSegWitWallet(ctx, store, walletName)
	case MultiSigType:
		ms, err := b.GetMultiSigWallet(ctx, store, walletName)
		if err != nil {
			return nil, err
		}

		return ms.wallet, nil
	default:
		return nil, errors.New(UnknownWalletTypeError)
	}
}
