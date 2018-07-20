package btc

import(
	"context"
	"errors"
	"time"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/salt"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

type credential struct {
	WalletName string
}

func pathCredentials(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "creds/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: "Wallet name",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathCredsRead,
		},

		HelpSynopsis:    pathCredsHelpSyn,
		HelpDescription: pathCredsHelpDesc,
	}
}

func (b *backend) pathCredsRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	walletName := d.Get("name").(string)
	if walletName == "" {
		return nil, errors.New("missing wallet name")
	}

	w, err := b.GetWallet(ctx, req.Storage, walletName)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return logical.ErrorResponse("Failed to create credentials for '" + walletName + "': wallet does not exist"), nil
	}

	cred := &credential{
		WalletName: walletName,
	}

	token, err := b.NewToken(ctx, req.Storage, cred, walletName)
	if err != nil {
		return nil, err
	}

	resp := b.Secret(SecretCredsType).Response(
		map[string]interface{}{"token": token},
		map[string]interface{}{"token": token},
	)

	resp.Secret.TTL = time.Duration(1 * time.Hour)
	resp.Secret.MaxTTL = time.Duration(1 * time.Hour)

	return resp, nil
}

func (b *backend) NewSaltedToken(ctx context.Context, s logical.Storage, config *salt.Config) (string, string, error) {
	token, err := uuid.GenerateUUID()
	if err != nil {
		return "", "", err
	}

	newSalt, err := salt.NewSalt(ctx, s, config)
	if err != nil {
		return "", "", err
	}

	return token, newSalt.SaltID(token), nil
}

func (b *backend) NewToken(ctx context.Context, store logical.Storage, cred *credential, walletName string) (string, error) {
	token, saltedToken, err := b.NewSaltedToken(ctx, store, nil)
	if err != nil {
		return "", err
	}

	entry, err := logical.StorageEntryJSON("creds/" + saltedToken, cred)
	if err != nil {
		return "", err
	}

	if err := store.Put(ctx, entry); err != nil {
		return "", err
	}

	return token, nil
}

func (b *backend) GetToken(ctx context.Context, s logical.Storage, token string) (*credential, error) {
	entry, err := s.Get(ctx, "creds/" + token)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var cred credential
	if err := entry.DecodeJSON(&cred); err != nil {
		return nil, err
	}

	return &cred, nil
}

const pathCredsHelpSyn = `
Creates access tokens for already generated wallet
`

const pathCredsHelpDesc = ``