package btc

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestTransaction(t *testing.T) {
	b, storage := getTestBackend(t)

	name := "test"
	network := "testnet"
	_, err := newWallet(t, b, storage, name, network)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := newAuthToken(t, b, storage, name)
	if err != nil {
		t.Fatal(err)
	}
	token := resp.Data["token"].(string)
	isMultisig := false
	rawTx := "0100000001be66e10da854e7aea9338c1f91cd489768d1d6d7189f586d7a3613f2a24d5396000000001976a914dd6cce9f255a8cc17bda8ba0373df8e861cb866e88acffffffff0123ce0100000000001976a9142bc89c2702e0e618db7d59eb5ce2f0f147b4075488ac0000000001000000"

	t.Run("Sign transaction for BIP44 wallet", func(t *testing.T) {

		resp, err := newSignature(t, b, storage, name, rawTx, isMultisig, token)
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("No response received")
		}
	})

	t.Run("Sign transaction with expired token should fail", func(t *testing.T) {
		t.Parallel()

		exp := InvalidTokenError
		_, err := newSignature(t, b, storage, name, rawTx, isMultisig, token)
		if err == nil {
			t.Fatal("Should have failed before")
		}
		if err.Error() != exp {
			t.Fatalf("Want: %v, got: %v", exp, err)
		}
	})

	t.Run("Sign transaction without token should fail", func(t *testing.T) {
		t.Parallel()

		token := ""
		exp := MissingTokenError
		_, err := newSignature(t, b, storage, name, rawTx, isMultisig, token)
		if err == nil {
			t.Fatal("Should have failed before")
		}
		if err.Error() != exp {
			t.Fatalf("Want: %v, got: %v", exp, err)
		}
	})

	t.Run("Sign transaction with invalid token should fail", func(t *testing.T) {
		t.Parallel()

		token := "testtoken"
		exp := InvalidTokenError
		_, err := newSignature(t, b, storage, name, rawTx, isMultisig, token)
		if err == nil {
			t.Fatal("Should have failed before")
		}
		if err.Error() != exp {
			t.Fatalf("Want: %v, got: %v", exp, err)
		}
	})
}

func newSignature(t *testing.T, b logical.Backend, store logical.Storage, name string, rawTx string, multisig bool, token string) (*logical.Response, error) {
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Path:      "transaction/" + name,
		Storage:   store,
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"multisig": multisig,
			"rawTx":    rawTx,
			"token":    token,
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
