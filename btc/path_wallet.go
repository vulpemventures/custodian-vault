package btc

import(
	"context"
	"errors"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/btcsuite/btcutil/hdkeychain"
)

type wallet struct {
	Network string
	MasterKey string
	DerivationPath []uint32
}

func pathWallet(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "wallet/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"network": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: "btc network type: mainnet | testnet",
			},
			"name": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: "wallet name",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathWalletRead,
			logical.UpdateOperation: b.pathWalletWrite,
		},

		HelpSynopsis:    pathWalletsHelpSyn,
		HelpDescription: pathWalletsHelpDesc,
	}
}

func (b *backend) pathWalletWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	network := d.Get("network").(string)
	if network == "" {
		return nil, errors.New("missing network")
	}

	walletName := d.Get("name").(string)
	if walletName == "" {
		return nil, errors.New("missing wallet name")
	}

	// return error if a wallet with same name has already been created
	w, err := b.GetWallet(ctx, req.Storage, walletName)
	if err != nil {
		return nil, err
	}
	if w != nil {
		return nil, errors.New("Wallet with name '" + walletName + "' already exists")
	}

	wallet, err := createWallet(network)
	if err != nil {
		return nil, err
	}
	
	// create storage entry
	entry, err := logical.StorageEntryJSON("wallet/" + walletName, wallet)
	if err != nil {
		return nil, err
	}

	// save in local storage
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathWalletRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	walletName := d.Get("name").(string)

	w, err := b.GetWallet(ctx, req.Storage, walletName)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return nil, nil
	}

	key, err := hdkeychain.NewKeyFromString(w.MasterKey)
	if err != nil {
		return nil, err
	}

	xprv, err := derivePrivKey(key, w.DerivationPath)
	if err != nil {
		return nil, err
	}
	
	xpub, err := derivePubKey(xprv)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"network": w.Network,
			"xpub": xpub.String(),
		},
	}, nil
}

func (b *backend) GetWallet(ctx context.Context, store logical.Storage, walletName string) (*wallet, error) {
	entry, err := store.Get(ctx, "wallet/" + walletName)
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

const pathWalletsHelpSyn = `
Creates a new wallet by specifying network and name
`

const pathWalletsHelpDesc = ``