// audit/storageaccounts/ensure_secure_transfer_required.go

package storageaccounts

import (
	"azure-security-hub/utils"
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
)

// EnsureSecureTransferRequired checks if secure transfer is required for all storage accounts
func EnsureSecureTransferRequired() string {
	// Load and print compliance info
	compliance, err := utils.LoadComplianceData("compliance/cis_microsoft_azure_foundations_benchmark_v3.0.0.json")
	if err != nil {
		log.Printf("└─[ERROR] Error loading compliance data: %v", err)
		return "NA"
	}
	utils.PrintComplianceInfo(compliance, "4.1")

	// Get Azure client
	client, err := utils.GetAzureClient()
	if err != nil {
		log.Printf("└─[ERROR] Failed to get Azure client: %v", err)
		return "NA"
	}

	log.Println("└─[*] Creating Storage Account client")
	storageClient, err := armstorage.NewAccountsClient(client.SubscriptionID, client.Cred, nil)
	if err != nil {
		log.Printf("└─[ERROR] Failed to create storage account client: %v", err)
		return "NA"
	}

	log.Println("└─[*] Listing storage accounts")
	pager := storageClient.NewListPager(nil)

	var nonCompliantAccounts []string
	accountsChecked := 0

	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			log.Printf("└─[ERROR] Failed to list storage accounts: %v", err)
			return "NA"
		}

		for _, account := range page.Value {
			accountsChecked++
			accountName := *account.Name

			log.Printf("  └─[*] Checking storage account: %s", accountName)

			if account.Properties == nil || !*account.Properties.EnableHTTPSTrafficOnly {
				nonCompliantAccounts = append(nonCompliantAccounts, accountName)
				log.Printf("    └─[FAIL] Secure transfer not required for: %s", accountName)
			} else {
				log.Printf("    └─[PASS] Secure transfer is required for: %s", accountName)
			}
		}
	}

	if accountsChecked == 0 {
		log.Println("└─[*] No storage accounts found")
		return "NA"
	}

	if len(nonCompliantAccounts) > 0 {
		log.Printf("└─[FAIL] Found %d storage accounts with secure transfer not required:", len(nonCompliantAccounts))
		for _, account := range nonCompliantAccounts {
			log.Printf("  └─[FAIL] %s", account)
		}
		return "FAIL"
	}

	log.Printf("└─[PASS] All %d storage accounts have secure transfer required enabled", accountsChecked)
	return "PASS"
}
