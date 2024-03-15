package utils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
	"github.com/joho/godotenv"
)

// LoginAzure performs Azure authentication and returns a subscription client
func LoginAzure() (*armsubscription.SubscriptionsClient, error) {
	// Load .env file
	envPath := filepath.Join("..", ".env")
	err := godotenv.Load(envPath)
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	// Get credentials from environment variables
	tenantID := os.Getenv("AZURE_TENANT_ID")
	clientID := os.Getenv("AZURE_CLIENT_ID")
	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")

	if tenantID == "" || clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("AZURE_TENANT_ID, AZURE_CLIENT_ID, or AZURE_CLIENT_SECRET is not set")
	}

	// Create credential object
	cred, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create credential: %v", err)
	}

	// Create subscription client
	client, err := armsubscription.NewSubscriptionsClient(cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription client: %v", err)
	}

	// Test the authentication by listing subscriptions
	pager := client.NewListPager(nil)
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to list subscriptions: %v", err)
		}
		for _, sub := range page.Value {
			fmt.Printf("Subscription: %s (%s)\n", *sub.DisplayName, *sub.SubscriptionID)
		}
	}

	fmt.Println("Azure authentication successful")
	return client, nil
}
