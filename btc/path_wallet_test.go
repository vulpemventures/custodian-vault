package btc

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestWallet(t *testing.T) {
	b, storage := getTestBackend(t)
	name := "test"
	network := "testnet"

	t.Run("Create BIP44 wallet", func(t *testing.T) {
		resp, err := newWallet(t, b, storage, name, network)
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("No response received")
		}
	})

	t.Run("Get BIP44 wallet info", func(t *testing.T) {
		t.Parallel()
		resp, err := getWallet(t, b, storage, name)
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("No response received")
		}
	})

	t.Run("Create BIP44 wallet without network should fail", func(t *testing.T) {
		t.Parallel()

		name := "testwallet"
		network := ""
		exp := MissingNetworkError
		_, err := newWallet(t, b, storage, name, network)
		if err == nil {
			t.Fatal("Should have failed before")
		}
		if err.Error() != exp {
			t.Fatalf("Want: %v, got: %v", exp, err)
		}
	})

	t.Run("Create an existing BIP44 wallet should fail", func(t *testing.T) {
		t.Parallel()

		exp := "Wallet with name '" + name + "' already exists"
		_, err := newWallet(t, b, storage, name, network)
		if err == nil {
			t.Fatal("Should have failed before")
		}
		if err.Error() != exp {
			t.Fatalf("Want: %v, got: %v", exp, err)
		}
	})
}

func newWallet(t *testing.T, b logical.Backend, store logical.Storage, name string, network string) (*logical.Response, error) {
	data := map[string]interface{}{"network": network}
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Storage:   store,
		Path:      "wallet/" + name,
		Operation: logical.UpdateOperation,
		Data:      data,
	})
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, resp.Error()
	}

	return resp, nil
}

func getWallet(t *testing.T, b logical.Backend, store logical.Storage, name string) (*logical.Response, error) {
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Storage:   store,
		Path:      "wallet/" + name,
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
