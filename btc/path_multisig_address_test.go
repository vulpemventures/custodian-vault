package btc

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestMultiSigAddress(t *testing.T) {
	b, storage := getTestBackend(t)

	exp := MissingTokenError
	_, err := b.HandleRequest(context.Background(), &logical.Request{
		Storage:   storage,
		Path:      "address/multisig/wallet1",
		Operation: logical.UpdateOperation,
	})
	if err == nil {
		t.Fatal("Should have failed before")
	}
	if err.Error() != exp {
		t.Fatalf("Want: %v, got: %v", exp, err)
	}

	exp = InvalidTokenError
	_, err = b.HandleRequest(context.Background(), &logical.Request{
		Storage:   storage,
		Path:      "address/multisig/wallet1",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"token": "testtoken",
		},
	})
	if err == nil {
		t.Fatal("Should have failed before")
	}
	if err.Error() != exp {
		t.Fatalf("Want: %v, got: %v", exp, err)
	}
}
