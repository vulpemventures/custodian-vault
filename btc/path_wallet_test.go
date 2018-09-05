package btc

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestWallet(t *testing.T) {
	b, storage := getTestBackend(t)

	t.Run("Create wallet", func(t *testing.T) {
		t.Parallel()

		exp := MissingNetworkError
		_, err := b.HandleRequest(context.Background(), &logical.Request{
			Storage:   storage,
			Path:      "wallet/testwallet",
			Operation: logical.UpdateOperation,
		})
		if err == nil {
			t.Fatal("Should have failed before")
		}
		if err.Error() != exp {
			t.Fatalf("Want: %v, got: %v", exp, err)
		}
	})

	t.Run("Get extendend public key", func(t *testing.T) {
		t.Parallel()

		resp, err := b.HandleRequest(context.Background(), &logical.Request{
			Storage:   storage,
			Path:      "wallet/wallet",
			Operation: logical.ReadOperation,
		})

		if err == nil && resp != nil {
			t.Fatal("Should have failed before")
		}
	})
}
