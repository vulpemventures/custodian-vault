package btc

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestSegWitWallet(t *testing.T) {
	b, storage := getTestBackend(t)
	name := "test"
	network := "testnet"

	t.Run("New BIP84 wallet", func(t *testing.T) {
		resp, err := newSegWitWallet(t, b, storage, name, network)
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("No response received")
		}

		t.Log("Mnemonic:", resp.Data["mnemonic"].(string))
	})

	t.Run("Get wallet info", func(t *testing.T) {
		t.Parallel()

		resp, err := getSegWitWallet(t, b, storage, name)
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("No response received")
		}

		t.Log("Wallet info:", resp.Data)
	})

	t.Run("New wallet without network should fail", func(t *testing.T) {
		t.Parallel()

		name := "testwallet"
		network := ""
		exp := MissingNetworkError
		_, err := newSegWitWallet(t, b, storage, name, network)
		if err == nil {
			t.Fatal("Should have failed before")
		}
		if err.Error() != exp {
			t.Fatalf("Want: %v, got: %v", exp, err)
		}
	})

	t.Run("New wallet with invalid network should fail", func(t *testing.T) {
		t.Parallel()

		name := "testwallet"
		network := "invaildnetwork"
		exp := InvalidNetworkError
		_, err := newSegWitWallet(t, b, storage, name, network)
		if err == nil {
			t.Fatal("Should have failed before")
		}
		if err.Error() != exp {
			t.Fatalf("Want: %v, got: %v", exp, err)
		}
	})

	t.Run("Create an existing wallet should fail", func(t *testing.T) {
		t.Parallel()

		exp := SegWitWalletAlreadyExistsError
		_, err := newSegWitWallet(t, b, storage, name, network)
		if err == nil {
			t.Fatal("Should have failed before")
		}
		if err.Error() != exp {
			t.Fatalf("Want: %v, got: %v", exp, err)
		}
	})
}

func newSegWitWallet(t *testing.T, b logical.Backend, store logical.Storage, name string, network string) (*logical.Response, error) {
	data := map[string]interface{}{"network": network}
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Storage:   store,
		Path:      PathSegWitWallet + name,
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

func getSegWitWallet(t *testing.T, b logical.Backend, store logical.Storage, name string) (*logical.Response, error) {
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Storage:   store,
		Path:      PathSegWitWallet + name,
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
