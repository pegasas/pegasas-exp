package kubernetes

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/patrickmn/go-cache"
	"k8s.io/client-go/tools/clientcmd"

	"k8s-delegater/pkg/common"
)

type AzureKubernetesConfig struct {
	credentials     *AzureKubernetesCredentials
	kubeConfigBytes []byte
	cache           *cache.Cache
}

func NewAzureKubernetesConfigFromEnv() *AzureKubernetesConfig {
	return &AzureKubernetesConfig{
		credentials:     NewAzureKubernetesCredentialsFromEnv(),
		kubeConfigBytes: nil,
		cache:           cache.New(5*time.Minute, 10*time.Minute),
	}
}

func NewAzureKubernetesConfig(credentials *AzureKubernetesCredentials) *AzureKubernetesConfig {
	return &AzureKubernetesConfig{
		credentials:     credentials,
		kubeConfigBytes: nil,
		cache:           cache.New(5*time.Minute, 10*time.Minute),
	}
}

func (aks *AzureKubernetesConfig) GetClientConfig(ctx context.Context) (config clientcmd.ClientConfig, err error) {
	kubeConfigBytes, err := aks.GetKubeConfigBytes(ctx)
	if err != nil {
		common.GetLog().Errorf("GetKubeConfigBytes error: %v", err)
		return nil, err
	}

	clientConfig, err := clientcmd.NewClientConfigFromBytes(kubeConfigBytes)
	if err != nil {
		common.GetLog().Errorf("clientcmd NewClientConfigFromBytes error: %v", err)
		return nil, err
	}

	rawConfig, err := clientConfig.RawConfig()
	if err != nil {
		common.GetLog().Errorf("clientConfig RawConfig error: %v", err)
		return nil, err
	}

	callback := func(t adal.Token) error {
		return nil
	}
	overrides := &clientcmd.ConfigOverrides{}
	for _, authInfo := range rawConfig.AuthInfos {
		serverID := authInfo.AuthProvider.Config["apiserver-id"]
		spt, err := adal.NewServicePrincipalTokenFromManagedIdentity(
			serverID,
			&adal.ManagedIdentityOptions{
				ClientID: aks.credentials.ClientID,
			},
			callback)
		if err != nil {
			common.GetLog().Errorf("adal NewServicePrincipalTokenFromManagedIdentity error: %v", err)
			return nil, err
		}
		spt.Refresh()
		token := spt.Token()
		log.Print(token.AccessToken)
		log.Print(token.ExpiresOn)
		authInfo.AuthProvider.Config["access-token"] = token.AccessToken
		authInfo.AuthProvider.Config["expires-on"] = string(token.ExpiresOn)
		overrides.AuthInfo.Token = token.AccessToken
	}

	clientConfig = clientcmd.NewDefaultClientConfig(rawConfig, overrides)
	return clientConfig, nil
}

func (aks *AzureKubernetesConfig) GetKubeConfigBytes(ctx context.Context) (kubeConfigBytes []byte, err error) {
	key := generateKey(aks.credentials.SubscriptionID, aks.credentials.ResourceGroup, aks.credentials.Cluster)
	value, found := aks.cache.Get(key)
	if found {
		common.GetLog().Info("get kube config bytes from cache")
		return value.([]byte), nil
	} else {
		value, _ := aks.GetKubeConfigBytesFromCluster(ctx)
		aks.cache.Set(key, value, cache.DefaultExpiration)
		return value, nil
	}
}

func (aks *AzureKubernetesConfig) GetKubeConfigBytesFromCluster(ctx context.Context) (kubeConfigBytes []byte, err error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		common.GetLog().Errorf("azidentity NewDefaultAzureCredential error: %v", err)
		return nil, err
	}

	managedClustersClient, err := armcontainerservice.NewManagedClustersClient(aks.credentials.SubscriptionID, cred, nil)
	if err != nil {
		common.GetLog().Errorf("armcontainerservice NewManagedClustersClient error: %v", err)
		return nil, err
	}

	accessProfile, err := managedClustersClient.ListClusterUserCredentials(ctx, aks.credentials.ResourceGroup, aks.credentials.Cluster, nil)
	if err != nil {
		common.GetLog().Errorf("managedClustersClient ListClusterUserCredentials error: %v", err)
		return nil, err
	}

	kubeConfigs := accessProfile.CredentialResults.Kubeconfigs
	if len(kubeConfigs) != 1 {
		common.GetLog().Errorf("accessProfile CredentialResults Kubeconfigs lens larger than 1")
		return nil, errors.New("accessProfile CredentialResults Kubeconfigs lens larger than 1")
	}
	return kubeConfigs[0].Value, nil
}

func generateKey(subscriptionID string, resourceGroup string, cluster string) string {
	var builder strings.Builder
	builder.WriteString(subscriptionID)
	builder.WriteString("_")
	builder.WriteString(resourceGroup)
	builder.WriteString("_")
	builder.WriteString(cluster)
	return builder.String()
}
