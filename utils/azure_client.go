package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
	"github.com/joho/godotenv"
)

type AzureClient struct {
	SubscriptionID string
	Cred           *azidentity.ClientSecretCredential
	SubClient      *armsubscription.SubscriptionsClient
}

// GetAzureClientFunc is a function type for getting an AzureClient
type GetAzureClientFunc func() (*AzureClient, error)

// GetAzureClient is a variable holding the function to get an AzureClient
var GetAzureClient GetAzureClientFunc = getAzureClient

// getAzureClient is the actual implementation
func getAzureClient() (*AzureClient, error) {
	// Load .env file from the current directory
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("[!] Warning: Error loading .env file: %v\n", err)
	}

	tenantID := os.Getenv("AZURE_TENANT_ID")
	clientID := os.Getenv("AZURE_CLIENT_ID")
	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")

	if tenantID == "" || clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("[-] AZURE_TENANT_ID, AZURE_CLIENT_ID, or AZURE_CLIENT_SECRET is not set")
	}

	// Debug output (be careful with clientSecret)
	// fmt.Printf("[*] Debug: AZURE_TENANT_ID=%s, AZURE_CLIENT_ID=%s, AZURE_CLIENT_SECRET=%s...\n",
	// 	tenantID, clientID, clientSecret[:5])

	cred, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
	if err != nil {
		return nil, fmt.Errorf("[-] Failed to create credential: %v", err)
	}

	subClient, err := armsubscription.NewSubscriptionsClient(cred, nil)
	if err != nil {
		return nil, fmt.Errorf("[-] Failed to create subscription client: %v", err)
	}

	var subscriptionID string
	pager := subClient.NewListPager(nil)
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			return nil, fmt.Errorf("[-] Failed to list subscriptions: %v", err)
		}
		for _, sub := range page.Value {
			subscriptionID = *sub.SubscriptionID
			fmt.Printf("[+] Found subscription: %s (%s)\n", *sub.DisplayName, *sub.SubscriptionID)
			break // Just use the first subscription
		}
	}

	if subscriptionID == "" {
		return nil, fmt.Errorf("[-] No subscription found")
	}

	fmt.Printf("[+] Using subscription: %s\n", subscriptionID)

	return &AzureClient{
		SubscriptionID: subscriptionID,
		Cred:           cred,
		SubClient:      subClient,
	}, nil
}

// TestAzureCredentials tests the Azure credentials
func TestAzureCredentials() error {
	fmt.Println("[*] Testing Azure credentials...")

	client, err := getAzureClient()
	if err != nil {
		return fmt.Errorf("[-] Failed to get Azure client: %v", err)
	}

	// Try to use the credential
	_, err = client.Cred.GetToken(context.Background(), policy.TokenRequestOptions{
		Scopes: []string{"https://management.azure.com/.default"},
	})
	if err != nil {
		return fmt.Errorf("[-] Failed to get token: %v", err)
	}

	fmt.Println("[+] Azure credentials test passed")
	return nil
}
