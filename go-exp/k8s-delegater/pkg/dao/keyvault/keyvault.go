package keyvault

import (
	"context"
)

type KeyVault interface {
	GetSecret(ctx context.Context, secretName string) (value string, err error)
}
