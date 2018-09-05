package btc

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestMultiSigPlugin(t *testing.T) {
	b, storage := getTestBackend(t)
	pubkeys := "04b353e7164238cc9106db773e0fa26507ede05ff78b31100ac47f9e328587bf2e7e740377ccbc29c5584bf13e2a05e47bfb8aa2529d4b2e91fe6879269da36c66,04179f874d9fc892f6b9d2cbc3aefea47f39f364d7a8e0c8bc3d6f1af1659160cb319f713260379bf26c0811505a471e9620fbb7404b4e8a0089a6d9cde301be3b"

	t.Run("Generate receiving address", func(t *testing.T) {
		var token string

		t.Run("Create auth token for wallet", func(t *testing.T) {
			t.Run("Path wallet/multisig/", func(t *testing.T) {
				t.Run("Create multisig wallet", func(t *testing.T) {
					resp, err := b.HandleRequest(context.Background(), &logical.Request{
						Storage:   storage,
						Path:      "wallet/multisig/testwallet",
						Operation: logical.UpdateOperation,
						Data: map[string]interface{}{
							"network": "testnet",
							"m":       2,
							"n":       3,
							"pubkeys": pubkeys,
						},
					})
					if err != nil || (resp != nil && resp.IsError()) {
						t.Fatalf("bad: resp: %#v\nerr: %v\n", resp, err)
					}
				})

				t.Run("Get multisig wallet info", func(t *testing.T) {
					resp, err := b.HandleRequest(context.Background(), &logical.Request{
						Storage:   storage,
						Path:      "wallet/multisig/testwallet",
						Operation: logical.ReadOperation,
					})

					if err != nil || (resp != nil && resp.IsError()) {
						t.Fatalf("bad: resp: %#v\nerr: %v\n", resp, err)
					}
				})
			})

			resp, err := b.HandleRequest(context.Background(), &logical.Request{
				Storage:   storage,
				Path:      "creds/multisig/testwallet",
				Operation: logical.ReadOperation,
			})
			if err != nil || (resp != nil && resp.IsError()) {
				t.Fatalf("bad: resp: %#v\nerr: %v\n", resp, err)
			}

			token = resp.Data["token"].(string)
		})

		resp, err := b.HandleRequest(context.Background(), &logical.Request{
			Storage:   storage,
			Path:      "address/multisig/testwallet",
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
