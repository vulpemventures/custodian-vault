package btc

import (
	"context"
	"encoding/hex"
	"errors"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathTransaction(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "transaction/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Wallet name",
			},
			"rawTx": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Raw transaction to be signed",
			},
			"mode": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Transaction type: standard | multisig | segwit ",
				Default:     "standard",
			},
			"token": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Auth token",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathTransactionWrite,
		},

		HelpSynopsis:    PathTransactionHelpSyn,
		HelpDescription: PathTransactionHelpDesc,
	}
}

func (b *backend) pathTransactionWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	walletName := d.Get("name").(string)
	if walletName == "" {
		return nil, errors.New(MissingWalletNameError)
	}
	mode := d.Get("mode").(string)
	walletType, err := parseMode(mode)
	if err != nil {
		return nil, err
	}

	t := d.Get("token").(string)
	if t == "" {
		return nil, errors.New(MissingTokenError)
	}

	// check if auth token is valid
	token, err := b.GetToken(ctx, req.Storage, t, walletType)
	if err != nil {
		return nil, err
	}
	if token == nil || token.WalletName != walletName {
		return nil, errors.New(InvalidTokenError)
	}

	w, err := getWalletByType(ctx, b, req.Storage, walletName, walletType)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return nil, errors.New(WalletNotFoundError)
	}

	rawTx := d.Get("rawTx").(string)
	if rawTx == "" {
		return nil, errors.New(MissingRawTxError)
	}

	seed := seedFromMnemonic(w.Mnemonic)

	masterKey, err := getMasterKey(seed, w.Network)
	if err != nil {
		return nil, err
	}

	storePath := PathAddress
	if walletType == SegWitType {
		storePath = PathSegWitAddress
	}
	storePath += walletName
	// derive key of last used address (for multisig is 0)
	childnum, err := b.GetLastUsedAddressIndex(ctx, req.Storage, storePath)
	if err != nil {
		return nil, err
	}

	derivedPrivKey, err := derivePrivKey(masterKey, append(w.DerivationPath, uint32(childnum)))
	privateKey, err := derivedPrivKey.ECPrivKey()
	if err != nil {
		return nil, err
	}

	// convert tx string to raw bytes
	rawTxBytes, err := hex.DecodeString(rawTx)
	if err != nil {
		return nil, err
	}

	// double sha256 before signing
	hashedRawTx := chainhash.DoubleHashB(rawTxBytes)
	signatureBytes, err := privateKey.Sign(hashedRawTx)
	if err != nil {
		return nil, err
	}

	// convert signature raw bytes to string
	signature := hex.EncodeToString(signatureBytes.Serialize())

	// revoke auth token
	err = b.RevokeToken(ctx, req.Storage, token, walletType)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"signature": signature,
		},
	}, nil
}

func parseMode(s string) (int, error) {
	switch s {
	case "standard":
		return StandardType, nil
	case "multisig":
		return MultiSigType, nil
	case "segwit":
		return SegWitType, nil
	default:
		return -1, errors.New(InvalidModeError)
	}
}
