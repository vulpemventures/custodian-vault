package btc

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

type multiSigWallet struct {
	*wallet
	M            int
	N            int
	RedeemScript string
	PublicKeys   []string
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
	w, err := b.GetMultiSigWallet(ctx, req.Storage, walletName)
	if err != nil {
		return nil, err
	}
	if w != nil {
		return nil, errors.New(MultiSigWalletAlreadyExistsError)
	}

	// create multisig wallet with params
	wallet, err := createMultiSigWallet(network, pubkeys, m, n)
	if err != nil {
		return nil, err
	}

	// create storage entry
	entry, err := logical.StorageEntryJSON(PathMultiSigWallet+walletName, wallet)
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
	entry, err := store.Get(ctx, PathMultiSigWallet+walletName)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var raw map[string]interface{}
	if err := entry.DecodeJSON(&raw); err != nil {
		return nil, err
	}

	return toMultiSig(raw)
}

func toMultiSig(raw map[string]interface{}) (*multiSigWallet, error) {
	// converting storage response to multiSig type sucks a lot!!
	// - DerivationPath is []interface{} so we need to scan the array and convert each element to uint32
	// 	- since each element is a json.Number, convertion is: json.Number -> int64 -> uint32
	// - M,N: json.Number -> int64 -> int
	// - PubKeys is also a []interface{} but []interface{} -> []string is quite smooth
	rawPath := raw["DerivationPath"].([]interface{})
	path := make([]uint32, len(rawPath))
	for i := range rawPath {
		t, err := rawPath[i].(json.Number).Int64()
		if err != nil {
			return nil, err
		}
		path[i] = uint32(t)
	}
	m, err := raw["M"].(json.Number).Int64()
	if err != nil {
		return nil, err
	}
	n, err := raw["N"].(json.Number).Int64()
	if err != nil {
		return nil, err
	}
	rawPubkeys := raw["PublicKeys"].([]interface{})
	pubkeys := make([]string, len(rawPubkeys))
	for i := range rawPubkeys {
		pubkeys[i] = rawPubkeys[i].(string)
	}

	mw := &multiSigWallet{
		RedeemScript: raw["RedeemScript"].(string),
		PublicKeys:   pubkeys,
		M:            int(m),
		N:            int(n),
	}
	mw.wallet = &wallet{
		Mnemonic:       raw["Mnemonic"].(string),
		Network:        raw["Network"].(string),
		DerivationPath: path,
	}

	return mw, nil
}
