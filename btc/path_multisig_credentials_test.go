package btc

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestMultiSigCredentials(t *testing.T) {
	b, storage := getTestBackend(t)

	exp := "Failed to create credentials for 'multisig_wallet1': wallet does not exist"
	_, err := b.HandleRequest(context.Background(), &logical.Request{
		Storage:   storage,
		Path:      "creds/multisig/wallet1",
		Operation: logical.ReadOperation,
	})
	if err == nil {
		t.Fatal("Should have failed before")
	}
	if err.Error() != exp {
		t.Fatalf("Want: %v, got: %v", exp, err)
	}
}
