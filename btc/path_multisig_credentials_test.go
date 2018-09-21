package btc

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestMultiSigCredentials(t *testing.T) {
	b, storage := getTestBackend(t)

	m := 2
	n := 3
	name := "test"
	network := "testnet"
	pubkeys := []string{"", ""}
	_, err := newMultiSigWallet(t, b, storage, name, network, m, n, pubkeys)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Create auth token for multisig wallet", func(t *testing.T) {
		t.Parallel()

		resp, err := newMultiSigAuthToken(t, b, storage, name)
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("No response received")
		}
	})

	t.Run("Create auth token for bad multisig wallet should fail", func(t *testing.T) {
		t.Parallel()

		name := "badwallet"
		exp := MultiSigWalletNotFoundError
		_, err := newMultiSigAuthToken(t, b, storage, name)
		if err == nil {
			t.Fatal("Should have failed before")
		}
		if err.Error() != exp {
			t.Fatalf("Want: %v, got: %v", exp, err)
		}
	})
}

func newMultiSigAuthToken(t *testing.T, b logical.Backend, store logical.Storage, name string) (*logical.Response, error) {
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Path:      "creds/multisig/" + name,
		Storage:   store,
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
