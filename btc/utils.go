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

func createWallet(network string, segwit bool) (*wallet, error) {
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

	// derivation path is:
	// * m/44'/0'/0'/0 for mainnet BIP44 wallet
	// * m/44'/1'/0'/0 for testnet BIP44 wallet
	// * m/49'/0'/0'/0 for mainnet BIP49 wallet
	// * m/49'/1'/0'/0 for testnet BIP49 wallet
	purpose := Purpose
	if segwit {
		purpose = SegwitPurpose
	}

	path := []uint32{purpose, CoinType[network], Account, Change}

	wallet := &wallet{
		Mnemonic:       mnemonic,
		Network:        network,
		DerivationPath: path,
		Segwit:         segwit,
	}

	return wallet, nil
}

func createMultiSigWallet(network string, pubkeys []string, m int, n int) (*multiSigWallet, error) {
	// first create a standard wallet
	w, err := createWallet(network, false)
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
		wallet:       w,
		PublicKeys:   pubkeys,
		RedeemScript: redeemScript,
		M:            m,
		N:            n,
	}

	return wallet, nil
}

func createSegWitWallet(network string) (*wallet, error) {
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

	path := []uint32{NativeSegwitPurpose, CoinType[network], Account, Change}

	wallet := &wallet{
		Mnemonic:       mnemonic,
		Network:        network,
		DerivationPath: path,
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
func deriveAddress(w *wallet, childnum int) (*address, error) {
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
	path := append(w.DerivationPath, uint32(childnum))
	key, err = derivePrivKey(key, path)
	if err != nil {
		return nil, err
	}

	// addr, err := getWalletAddress(key, net, w.Segwit)
	addr := ""
	if w.Segwit {
		addr, err = bip49Address(key, net)
		if err != nil {
			return nil, err
		}
	} else {
		// generate new address for derived key
		rawAddress, err := key.Address(net)
		if err != nil {
			return nil, err
		}
		addr = rawAddress.String()
	}

	return &address{
		Childnum:    childnum,
		LastAddress: addr,
	}, nil
}

func deriveSegWitAddress(w *wallet, childnum int) (*address, error) {
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
	path := append(w.DerivationPath, uint32(childnum))
	key, err = derivePrivKey(key, path)
	if err != nil {
		return nil, err
	}

	pubkey, err := key.ECPubKey()
	if err != nil {
		return nil, err
	}

	witnessProgram := btcutil.Hash160(pubkey.SerializeCompressed())
	// version 0 of P2wPKH requires first 20 bytes of witness program
	addr, err := btcutil.NewAddressWitnessPubKeyHash(witnessProgram[:20], net)
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
	case RegTest:
		return &chaincfg.RegressionNetParams, nil
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

func bip49Address(key *hdkeychain.ExtendedKey, net *chaincfg.Params) (string, error) {
	pubkey, err := key.ECPubKey()
	if err != nil {
		return "", err
	}
	keyHash := btcutil.Hash160(pubkey.SerializeCompressed())
	// scriptSig for P2SHP2WPKH (version 0) is <0 <20-byte-key-hash>> as stated in BIP141
	scriptSig, err := txscript.NewScriptBuilder().AddOp(txscript.OP_0).AddData(keyHash[:20]).Script()
	if err != nil {
		return "", err
	}
	rawAddress, err := btcutil.NewAddressScriptHash(scriptSig, net)
	if err != nil {
		return "", err
	}

	return rawAddress.String(), nil
}
