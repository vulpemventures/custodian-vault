package btc

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestSegWitAddress(t *testing.T) {
	b, storage := getTestBackend(t)

	name := "test"
	network := "testnet"
	_, err := newSegWitWallet(t, b, storage, name, network)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := newSegWitAuthToken(t, b, storage, name)
	if err != nil {
		t.Fatal(err)
	}
	token := resp.Data["token"].(string)

	t.Run("Get address for native segwit wallet", func(t *testing.T) {
		resp, err := newSegWitAddress(t, b, storage, name, token)
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("No response received")
		}

		address := resp.Data["address"].(string)
		if !strings.HasPrefix(address, "tb") {
			t.Fatal("Invalid address:", address)
		}
		t.Log("Address:", address)
	})

	t.Run("Get address with expired auth token should fail", func(t *testing.T) {
		t.Parallel()

		exp := InvalidTokenError
		_, err := newSegWitAddress(t, b, storage, name, token)
		if err == nil {
			t.Fatal("Should have failed before")
		}
		if err.Error() != exp {
			t.Fatalf("Want: %v, got %v", exp, err)
		}
	})

	t.Run("Get address without auth token should fail", func(t *testing.T) {
		t.Parallel()

		token := ""
		exp := MissingTokenError
		_, err := newSegWitAddress(t, b, storage, name, token)
		if err == nil {
			t.Fatal("Should have failed before")
		}
		if err.Error() != exp {
			t.Fatalf("Want: %v, got: %v", exp, err)
		}
	})

	t.Run("Get address with invalid auth token should fail", func(t *testing.T) {
		t.Parallel()

		token := "testtoken"
		exp := InvalidTokenError
		_, err := newSegWitAddress(t, b, storage, name, token)
		if err == nil {
			t.Fatal("Should have failed before")
		}
		if err.Error() != exp {
			t.Fatalf("Want: %v, got: %v", exp, err)
		}
	})
}

func newSegWitAddress(t *testing.T, b logical.Backend, store logical.Storage, name string, token string) (*logical.Response, error) {
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Storage:   store,
		Path:      "address/segwit/" + name,
		Operation: logical.UpdateOperation,
		Data:      map[string]interface{}{"token": token},
	})
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, resp.Error()
	}

	return resp, nil
}
