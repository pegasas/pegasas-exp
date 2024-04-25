package keyvault

import (
	"context"
)

type LocalKeyVaultClient struct {
	secret string
}

func NewLocalKeyVaultClient(secret string) (*LocalKeyVaultClient, error) {
	return &LocalKeyVaultClient{secret: secret}, nil
}

func (kv *LocalKeyVaultClient) GetSecret(ctx context.Context, secretName string) (value string, err error) {
	return kv.secret, nil
}
