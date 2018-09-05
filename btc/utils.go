package btc

import (
	"bytes"
	"encoding/hex"
	"errors"

	"github.com/tyler-smith/go-bip39"

	"github.com/btcsuite/btcutil"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil/hdkeychain"
)

func seedFromMnemonic(mnemonic string) []byte {
	return bip39.NewSeed(mnemonic, "")
}

func createWallet(network string) (*wallet, error) {
	// generate entropy for mnemonic
	entropy, err := bip39.NewEntropy(EntropyBitSize)
	if err != nil {
		return nil, err
	}
	// generate mnemonic
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, err
	}

	// derivation path hard coded to m/44'/0'/0'/0 for mainnet
	// and m/44'/1'/0'/0 for testnet as specified in BIP44
	// support for BIP49 will be added later
	// hardened key derivation starts at index 2147483648 (hex 0x80000000)
	var path []uint32
	if network == MainNet {
		path = []uint32{Purpose, CoinTypeMainNet, Account, Change}
	}
	if network == TestNet {
		path = []uint32{Purpose, CoinTypeTestNet, Account, Change}
	}
	wallet := &wallet{
		Mnemonic:       mnemonic,
		Network:        network,
		DerivationPath: path,
	}

	return wallet, nil
}

func createMultiSigWallet(network string, pubkeys []string, m int, n int) (*multiSigWallet, error) {
	// first create a standard wallet
	w, err := createWallet(network)
	if err != nil {
		return nil, err
	}

	seed := seedFromMnemonic(w.Mnemonic)

	// get master key
	masterKey, err := getMasterKey(seed, w.Network)
	if err != nil {
		return nil, err
	}

	// derive private key at path m/44'/1'/0'/0/0
	// this is the private key used to sign raw trasactions
	privkey, err := derivePrivKey(masterKey, append(w.DerivationPath, MultiSigDefaultAddressIndex))
	if err != nil {
		return nil, err
	}

	// get public key from derived key (not in extended key format)
	pubkey, err := privkey.ECPubKey()
	if err != nil {
		return nil, err
	}
	net, err := getNetworkFromString(w.Network)
	if err != nil {
		return nil, err
	}
	pubkeybytes, err := btcutil.NewAddressPubKey(pubkey.SerializeUncompressed(), net)
	if err != nil {
		return nil, err
	}

	// add public key to list of public keys to build multisig wallet
	pubkeys = append(pubkeys, pubkeybytes.String())

	// create the redemption script
	redeemScript, err := newRedeemScript(m, n, pubkeys)
	if err != nil {
		return nil, err
	}

	wallet := &multiSigWallet{
		Mnemonic:       w.Mnemonic,
		Network:        w.Network,
		DerivationPath: w.DerivationPath,
		PublicKeys:     pubkeys,
		RedeemScript:   redeemScript,
		M:              m,
		N:              n,
	}

	return wallet, nil
}

// returns derived private key at `path` for `key` in extended key format
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

// returns public key for derived `key` in extended key format
func derivePubKey(key *hdkeychain.ExtendedKey) (*hdkeychain.ExtendedKey, error) {
	// neuter private key to get public key
	pubKey, err := key.Neuter()
	if err != nil {
		return nil, err
	}

	return pubKey, nil
}

// derive address for wallet `w` at path m/44'/0'/0'/0/childnum
func deriveAddress(w *wallet, childnum uint32) (*address, error) {
	net, err := getNetworkFromString(w.Network)
	if err != nil {
		return nil, err
	}

	seed := seedFromMnemonic(w.Mnemonic)

	// master key
	key, err := hdkeychain.NewMaster(seed, net)
	if err != nil {
		return nil, err
	}

	// append childnum to derivation path and derive private key
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
		Childnum:    childnum,
		LastAddress: addr.String(),
	}, nil
}

func getNetworkFromString(network string) (*chaincfg.Params, error) {
	switch network {
	case MainNet:
		return &chaincfg.MainNetParams, nil
	case TestNet:
		return &chaincfg.TestNet3Params, nil
	default:
		return nil, errors.New(InvalidNetworkError)
	}
}

// returns master private key from seed
func getMasterKey(seed []byte, network string) (*hdkeychain.ExtendedKey, error) {
	net, err := getNetworkFromString(network)
	if err != nil {
		return nil, err
	}

	key, err := hdkeychain.NewMaster(seed, net)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func newRedeemScript(m int, n int, pubkeys []string) (string, error) {
	// check if params are valid:
	// 1 <= n <= 7 && 1 <= m <= n
	if n < MinMultiSigN || n > MaxMultiSigN {
		return "", errors.New(NOutOfRangeError)
	}
	if m < MinMultiSigN || m > n {
		return "", errors.New(MOutOfRangeError)
	}

	if len(pubkeys) != n {
		return "", errors.New("Invalid number of pub keys: provided " + string(len(pubkeys)) + " expected " + string(n))
	}

	// get OP Code for m and n
	mOPCode := txscript.OP_1 + (m - 1)
	nOPCode := txscript.OP_1 + (n - 1)
	// multisig redeemScript format:
	// <OP_m> <A pubkey> <B pubkey> <C pubkey>... <OP_n> OP_CHECKMULTISIG
	var redeemScript bytes.Buffer
	redeemScript.WriteByte(byte(mOPCode))
	for _, pubkey := range pubkeys {
		pubkeybytes, err := hex.DecodeString(pubkey)
		if err != nil {
			return "", err
		}
		redeemScript.WriteByte(byte(len(pubkeybytes))) //PUSH
		redeemScript.Write(pubkeybytes)                //<pubkey>
	}
	redeemScript.WriteByte(byte(nOPCode))
	redeemScript.WriteByte(byte(txscript.OP_CHECKMULTISIG))

	return hex.EncodeToString(redeemScript.Bytes()), nil
}

// returns p2sh address derived from redeem script of multisig
func getMultiSigAddress(redeemScript string, network string) (string, error) {
	net, err := getNetworkFromString(network)
	if err != nil {
		return "", err
	}

	redeemScriptBytes, err := hex.DecodeString(redeemScript)
	if err != nil {
		return "", err
	}

	address, err := btcutil.NewAddressScriptHash(redeemScriptBytes, net)
	if err != nil {
		return "", nil
	}

	return address.String(), nil
}
