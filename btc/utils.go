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

	hkStart := uint32(0x80000000)
	wallet := &wallet{
		Seed: seed,
		Network: network,
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
	net, err := getNetworkFromString(w.Network)
	if err != nil {
		return nil, err
	}

	// get path master key
	key, err := hdkeychain.NewMaster(w.Seed, net)
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
	addr, err := key.Address(net)
	if err != nil {
		return nil, err
	}

	return &address{
		Childnum: childnum,
		LastAddress: addr.String(),
	}, nil
}

func getNetworkFromString(network string) (*chaincfg.Params, error) {
	switch network {
	case "mainnet":
		return &chaincfg.MainNetParams, nil
	case "testnet":
		return &chaincfg.TestNet3Params, nil
	default:
			return nil, errors.New("Invalid network")
	}
}