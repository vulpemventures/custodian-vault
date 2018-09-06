package btc

import (
	"testing"
)

func TestMultiSigPlugin(t *testing.T) {
	b, storage := getTestBackend(t)
	name := "test"
	network := "testnet"
	m := 2
	n := 3
	pubkeys := []string{
		"04b353e7164238cc9106db773e0fa26507ede05ff78b31100ac47f9e328587bf2e7e740377ccbc29c5584bf13e2a05e47bfb8aa2529d4b2e91fe6879269da36c66",
		"04179f874d9fc892f6b9d2cbc3aefea47f39f364d7a8e0c8bc3d6f1af1659160cb319f713260379bf26c0811505a471e9620fbb7404b4e8a0089a6d9cde301be3b",
	}
	rawTx := "0100000001be66e10da854e7aea9338c1f91cd489768d1d6d7189f586d7a3613f2a24d5396000000001976a914dd6cce9f255a8cc17bda8ba0373df8e861cb866e88acffffffff0123ce0100000000001976a9142bc89c2702e0e618db7d59eb5ce2f0f147b4075488ac0000000001000000"

	t.Logf("Test Parameters\nWallet name: %s\nWallet network: %s\nM: %d\nN: %d\nPublic keys: %s\nRaw transaction: %s\n", name, network, m, n, pubkeys, rawTx)
	// create new multisig wallet
	_, err := newMultiSigWallet(t, b, storage, name, network, m, n, pubkeys)
	if err != nil {
		t.Fatal(err)
	}

	// retrieve wallet info
	resp, err := getMultiSigWallet(t, b, storage, name)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v", resp.Data)

	// create auth token for multisig address
	resp, err = newMultiSigAuthToken(t, b, storage, name)
	if err != nil {
		t.Fatal(err)
	}
	token := resp.Data["token"].(string)

	// derive address
	resp, err = newMultiSigAddress(t, b, storage, name, token)
	if err != nil {
		t.Fatal(err)
	}
	firstAddress := resp.Data["address"].(string)

	// create auth token for multisig address
	resp, err = newMultiSigAuthToken(t, b, storage, name)
	if err != nil {
		t.Fatal(err)
	}
	token = resp.Data["token"].(string)

	// derive address
	resp, err = newMultiSigAddress(t, b, storage, name, token)
	if err != nil {
		t.Fatal(err)
	}
	secondAddress := resp.Data["address"].(string)

	// check that they don't match
	if firstAddress != secondAddress {
		t.Fatalf("Different requests generated the different addresses: %s, %s", firstAddress, secondAddress)
	}
	t.Logf("Multisig wallet address: %s", firstAddress)

	// create auth token for signature
	resp, err = newMultiSigAuthToken(t, b, storage, name)
	if err != nil {
		t.Fatal(err)
	}
	token = resp.Data["token"].(string)

	// create signature for raw transaction
	resp, err = newSignature(t, b, storage, name, rawTx, true, token)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Partial signature for raw tx: %s", resp.Data["signature"].(string))
}
