package btc

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestCredentials(t *testing.T) {
	b, storage := getTestBackend(t)

	name := "test"
	network := "testnet"
	_, err := newWallet(t, b, storage, name, network, !segwitCompatible)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Create auth token for wallet", func(t *testing.T) {
		t.Parallel()

		resp, err := newAuthToken(t, b, storage, name)
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("No response received")
		}
	})

	t.Run("Create auth token for bad wallet should fail", func(t *testing.T) {
		t.Parallel()

		name := "badwallet"
		exp := WalletNotFoundError
		_, err := newAuthToken(t, b, storage, name)
		if err == nil {
			t.Fatal("Should have failed before")
		}
		if err.Error() != exp {
			t.Fatalf("Want: %v, got: %v", exp, err)
		}
	})
}

func newAuthToken(t *testing.T, b logical.Backend, store logical.Storage, name string) (*logical.Response, error) {
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Path:      "creds/" + name,
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
