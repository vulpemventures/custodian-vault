package btc

import (
	"testing"
)

func TestPlugin(t *testing.T) {
	b, storage := getTestBackend(t)
	name := "test"
	network := "testnet"
	rawTx := "0100000001be66e10da854e7aea9338c1f91cd489768d1d6d7189f586d7a3613f2a24d5396000000001976a914dd6cce9f255a8cc17bda8ba0373df8e861cb866e88acffffffff0123ce0100000000001976a9142bc89c2702e0e618db7d59eb5ce2f0f147b4075488ac0000000001000000"

	t.Logf("Test Parameters\nWallet name: %s\nWallet network: %s\nRaw transaction: %s\n", name, network, rawTx)
	// create new BIP44 wallet
	resp, err := newWallet(t, b, storage, name, network)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Mnemonic: %s", resp.Data["mnemonic"].(string))

	// retrieve wallet info
	resp, err = getWallet(t, b, storage, name)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Tpub: %s", resp.Data["xpub"].(string))

	// create auth token for first address
	resp, err = newAuthToken(t, b, storage, name)
	if err != nil {
		t.Fatal(err)
	}
	token := resp.Data["token"].(string)

	// derive first address
	resp, err = newAddress(t, b, storage, name, token)
	if err != nil {
		t.Fatal(err)
	}
	firstAddress := resp.Data["address"].(string)

	// create auth token for second address
	resp, err = newAuthToken(t, b, storage, name)
	if err != nil {
		t.Fatal(err)
	}
	token = resp.Data["token"].(string)

	// derive second address
	resp, err = newAddress(t, b, storage, name, token)
	if err != nil {
		t.Fatal(err)
	}
	secondAddress := resp.Data["address"].(string)

	// check that they don't match
	if firstAddress == secondAddress {
		t.Fatal("Different requests generated the same address: " + firstAddress)
	}
	t.Logf("Derived addresses: %s, %s", firstAddress, secondAddress)

	// create auth token for signature
	resp, err = newAuthToken(t, b, storage, name)
	if err != nil {
		t.Fatal(err)
	}
	token = resp.Data["token"].(string)

	// create signature for raw transaction
	resp, err = newSignature(t, b, storage, name, rawTx, false, token)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Signature for raw tx: %s", resp.Data["signature"].(string))
}
