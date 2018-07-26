package btc

import (
	"testing"
	"context"

	"github.com/hashicorp/vault/logical"
)

func TestPlugin(t *testing.T) {
	b, storage := getTestBackend(t)

	t.Run("Generate receiving address", func(t *testing.T) {
		var token string

		t.Run("Create auth token for wallet", func(t *testing.T) {
			t.Run("Path wallet/", func(t *testing.T) {
				t.Run("Create wallet", func(t *testing.T) {
					resp, err := b.HandleRequest(context.Background(), &logical.Request{
						Storage: storage,
						Path: "wallet/testwallet",
						Operation: logical.UpdateOperation,
						Data: map[string]interface{}{
							"network": "testnet",
						},
					})
					if err != nil || (resp != nil && resp.IsError()) {
						t.Fatalf("bad: resp: %#v\nerr: %v\n", resp, err)
					}
				})

				t.Run("Get extended public key", func(t *testing.T) {
					resp, err := b.HandleRequest(context.Background(), &logical.Request{
						Storage: storage,
						Path: "wallet/testwallet",
						Operation: logical.ReadOperation,
					})
			
					if err != nil || (resp != nil && resp.IsError()) {
						t.Fatalf("bad: resp: %#v\nerr: %v\n", resp, err)
					}
				})
			})

			resp, err := b.HandleRequest(context.Background(), &logical.Request{
				Storage: storage,
				Path: "creds/testwallet",
				Operation: logical.ReadOperation,
			})
			if err != nil || (resp != nil && resp.IsError()) {
				t.Fatalf("bad: resp: %#v\nerr: %v\n", resp, err)
			}
		
			token = resp.Data["token"].(string)
		})

		resp, err := b.HandleRequest(context.Background(), &logical.Request{
			Storage: storage,
			Path: "address/testwallet",
			Operation: logical.UpdateOperation,
			Data: map[string]interface{}{
				"token": token,
			},
		})
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: resp: %#v\nerr: %v\n", resp, err)
		}
	})
}