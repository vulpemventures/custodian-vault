package btc

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestMultiSigAddress(t *testing.T) {
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

	resp, err := newMultiSigAuthToken(t, b, storage, name)
	token := resp.Data["token"].(string)

	t.Run("Get address for multisig wallet", func(t *testing.T) {
		resp, err := newMultiSigAddress(t, b, storage, name, token)
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("No response received")
		}
	})

	t.Run("Get address with expired auth token should fail", func(t *testing.T) {
		t.Parallel()

		exp := InvalidTokenError
		_, err := newMultiSigAddress(t, b, storage, name, token)
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
		_, err := newMultiSigAddress(t, b, storage, name, token)
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
		_, err := newMultiSigAddress(t, b, storage, name, token)
		if err == nil {
			t.Fatal("Should have failed before")
		}
		if err.Error() != exp {
			t.Fatalf("Want: %v, got: %v", exp, err)
		}
	})
}

func newMultiSigAddress(t *testing.T, b logical.Backend, store logical.Storage, name string, token string) (*logical.Response, error) {
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Storage:   store,
		Path:      "address/multisig/" + name,
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
