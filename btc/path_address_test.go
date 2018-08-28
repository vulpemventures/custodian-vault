package btc

import (
	"testing"
	"context"

	"github.com/hashicorp/vault/logical"
)

func TestAddress(t *testing.T) {
	b, storage := getTestBackend(t)

	exp := "missing auth token"
	_, err := b.HandleRequest(context.Background(), &logical.Request{
		Storage: storage,
		Path: "address/wallet1",
		Operation: logical.UpdateOperation,
	})
	if err == nil {
		t.Fatal("Should have failed before")
	}
	if err.Error() != exp {
		t.Fatalf("Want: %v, got: %v", exp, err)
	}

	exp = "token not found"
	_, err = b.HandleRequest(context.Background(), &logical.Request{
		Storage: storage,
		Path: "address/wallet1",
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