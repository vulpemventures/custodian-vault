package btc

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestMultiSigWallet(t *testing.T) {
	b, storage := getTestBackend(t)
	m := 2
	n := 3
	name := "test"
	network := "testnet"
	pubkeys := []string{
		"04b353e7164238cc9106db773e0fa26507ede05ff78b31100ac47f9e328587bf2e7e740377ccbc29c5584bf13e2a05e47bfb8aa2529d4b2e91fe6879269da36c66",
		"04179f874d9fc892f6b9d2cbc3aefea47f39f364d7a8e0c8bc3d6f1af1659160cb319f713260379bf26c0811505a471e9620fbb7404b4e8a0089a6d9cde301be3b",
	}

	t.Run("Create multisig wallet", func(t *testing.T) {
		// not checking response cause this returns nil if successful
		_, err := newMultiSigWallet(t, b, storage, name, network, m, n, pubkeys)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Get multisig wallet info", func(t *testing.T) {
		t.Parallel()

		resp, err := getMultiSigWallet(t, b, storage, name)
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("No response received")
		}
		t.Log(resp.Data)
	})

	t.Run("Create multisig wallet with missing params should fail", func(t *testing.T) {
		t.Parallel()

		err := badRequest(t, b, storage, name, "", m, n, pubkeys)
		if err != nil {
			t.Fatal(err)
		}

		err = badRequest(t, b, storage, name, network, 0, n, pubkeys)
		if err != nil {
			t.Fatal(err)
		}

		err = badRequest(t, b, storage, name, network, m, 0, pubkeys)
		if err != nil {
			t.Fatal(err)
		}

		err = badRequest(t, b, storage, name, network, m, n, []string{})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Create an existing multisig wallet should fail", func(t *testing.T) {
		t.Parallel()

		exp := MultiSigWalletAlreadyExistsError
		_, err := newMultiSigWallet(t, b, storage, name, network, m, n, pubkeys)
		if err == nil {
			t.Fatal("Should have failed before")
		}
		if err.Error() != exp {
			t.Fatalf("Want: %v, got: %v", exp, err)
		}
	})
}

func newMultiSigWallet(t *testing.T, b logical.Backend, store logical.Storage, name string, network string, m int, n int, pubkeys []string) (*logical.Response, error) {
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Storage:   store,
		Path:      "wallet/multisig/" + name,
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"m":       m,
			"n":       n,
			"pubkeys": pubkeys,
			"network": network,
		},
	})
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, resp.Error()
	}

	return resp, nil
}

func getMultiSigWallet(t *testing.T, b logical.Backend, store logical.Storage, name string) (*logical.Response, error) {
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Storage:   store,
		Path:      "wallet/multisig/" + name,
		Operation: logical.ReadOperation,
	})
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, resp.Error()
	}

	return resp, nil
}

func badRequest(t *testing.T, b logical.Backend, store logical.Storage, name string, network string, m int, n int, pubkeys []string) error {
	exp := MissingNetworkError
	if m == 0 {
		exp = InvalidMError
	}
	if n == 0 {
		exp = InvalidNError
	}
	if len(pubkeys) == 0 {
		exp = MissingPubKeysError
	}

	_, err := newMultiSigWallet(t, b, store, name, network, m, n, pubkeys)
	if err == nil {
		return errors.New("Should have failed before")
	}
	if err.Error() != exp {
		return errors.New("Want: " + exp + ", got: " + err.Error())
	}

	return nil
}
