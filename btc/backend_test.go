package btc

import (
	"testing"
	"context"

	"github.com/hashicorp/vault/logical"
)

func getTestBackend(t * testing.T) (logical.Backend, logical.Storage) {
	storage := &logical.InmemStorage{}

	config := logical.TestBackendConfig()
	config.StorageView = storage

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	if b == nil {
		t.Fatalf("Unable to create backend")
	}

	return b, storage
}