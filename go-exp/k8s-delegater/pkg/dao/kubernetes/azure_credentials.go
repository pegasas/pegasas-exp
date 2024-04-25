package kubernetes

import (
	"os"
)

const (
	subscriptionId = "SUBSCRIPTION_ID"
	resourceGroup  = "RESOURCE_GROUP"
	cluster        = "CLUSTER"
	clientId       = "CLIENT_ID"
)

type AzureKubernetesCredentials struct {
	SubscriptionID string `json:"subscription-id"`
	ResourceGroup  string `json:"resource-group"`
	Cluster        string `json:"cluster"`
	ClientID       string `json:"client-id"`
}

func NewAzureKubernetesCredentialsFromEnv() *AzureKubernetesCredentials {
	subscriptionID := os.Getenv(subscriptionId)
	resourceGroup := os.Getenv(resourceGroup)
	cluster := os.Getenv(cluster)
	clientID := os.Getenv(clientId)
	return &AzureKubernetesCredentials{
		SubscriptionID: subscriptionID,
		ResourceGroup:  resourceGroup,
		Cluster:        cluster,
		ClientID:       clientID,
	}
}
