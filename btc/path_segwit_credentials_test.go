package btc

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestSegWitCredentials(t *testing.T) {
	b, storage := getTestBackend(t)
	name := "test"
	network := "testnet"
	_, err := newSegWitWallet(t, b, storage, name, network)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Create auth token for native segwit wallet", func(t *testing.T) {
		t.Parallel()

		resp, err := newSegWitAuthToken(t, b, storage, name)
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("No response received")
		}
	})

	t.Run("Create auth token for bad native segwit wallet should fail", func(t *testing.T) {
		t.Parallel()

		name := "badwallet"
		exp := SegWitWalletNotFoundError
		_, err := newSegWitAuthToken(t, b, storage, name)
		if err == nil {
			t.Fatal("Should have failed before")
		}
		if err.Error() != exp {
			t.Fatalf("Want: %v, got: %v", exp, err)
		}
	})
}

func newSegWitAuthToken(t *testing.T, b logical.Backend, store logical.Storage, name string) (*logical.Response, error) {
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Path:      "creds/segwit/" + name,
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
