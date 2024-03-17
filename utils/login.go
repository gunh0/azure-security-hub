// utils/login.go
package utils

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
)

func LoginAzure() (*armsubscription.SubscriptionsClient, error) {
	client, err := GetAzureClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get Azure client: %v", err)
	}
	fmt.Println("Azure authentication successful")
	return client.SubClient, nil
}
