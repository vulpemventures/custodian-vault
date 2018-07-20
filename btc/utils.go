package btc

import(
	"errors"

	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
)

func createWallet(network string) (*wallet, error) {
	// generate a random seed
	seed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
	if err!= nil {
		return nil, err
	}

	var net *chaincfg.Params
	switch network {
	case "mainnet":
		net = &chaincfg.MainNetParams
	case "testnet":
		net = &chaincfg.TestNet3Params
	default:
			return nil, errors.New("Invalid network")
	}

	// generate master key 
	masterKey, err := hdkeychain.NewMaster(seed, net)
	hkStart := uint32(0x80000000)
	wallet := &wallet{
		Network: network,
		MasterKey: masterKey.String(),
		DerivationPath: []uint32{hkStart + 44, hkStart, hkStart}, // hardcoded m/44'/0'/0'
	}

	return wallet, nil
}

func derivePrivKey(key *hdkeychain.ExtendedKey, path []uint32) (*hdkeychain.ExtendedKey, error) {
	derivedKey := key

	for _, childNum := range path {
		var err error
		derivedKey, err = derivedKey.Child(childNum)
		if err != nil {
			return nil, err
		}
	}

	return derivedKey, nil
}

func derivePubKey(key *hdkeychain.ExtendedKey) (*hdkeychain.ExtendedKey, error) {
	// neuter private key to get public key
	pubKey, err := key.Neuter()
	if err != nil {
		return nil, err
	}

	return pubKey, nil
}

func deriveAddress(w *wallet, childnum uint32) (*address, error) {
	// get path master key
	key, err := hdkeychain.NewKeyFromString(w.MasterKey)
	if err != nil {
		return nil, err
	}

	// append childnum to derivation path
	path := append(w.DerivationPath, childnum)
	key, err = derivePrivKey(key, path)
	if err != nil {
		return nil, err
	}

	// generate new address for derived key
	net := getNetworkFromString(w.Network)
	addr, err := key.Address(net)
	if err != nil {
		return nil, err
	}

	return &address{
		Childnum: childnum,
		LastAddress: addr.String(),
	}, nil
}

func getNetworkFromString(network string) (*chaincfg.Params) {
	if network == "testnet" {
		return &chaincfg.TestNet3Params
	}

	return &chaincfg.MainNetParams
}