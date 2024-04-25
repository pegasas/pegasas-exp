package keyvault

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"

	"k8s-delegater/pkg/common"
)

type KeyVaultClient struct {
	client *azsecrets.Client
}

func NewKeyVaultClient() (*KeyVaultClient, error) {
	keyVaultName := common.GetKeyVaultName()
	keyVaultUrl := fmt.Sprintf("https://%s.vault.azure.net/", keyVaultName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		common.GetLog().Errorf("azidentity NewDefaultAzureCredential error: %v", err)
		return nil, err
	}

	client, err := azsecrets.NewClient(keyVaultUrl, cred, nil)
	if err != nil {
		common.GetLog().Errorf("azsecrets NewClient error: %v", err)
		return nil, err
	}

	return &KeyVaultClient{client: client}, nil
}

func (kv *KeyVaultClient) GetSecret(ctx context.Context, secretName string) (value string, err error) {
	resp, err := kv.client.GetSecret(ctx, secretName, "", nil)
	if err != nil {
		common.GetLog().Errorf("client GetSecret error: %v", err)
		return "", err
	}

	if resp.Value == nil {
		common.GetLog().Errorf("SecretBundle Value null")
		return "", errors.New("SecretBundle Value null")
	}

	return *(resp.Value), nil
}
