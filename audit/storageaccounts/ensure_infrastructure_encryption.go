// audit/storageaccounts/ensure_infrastructure_encryption.go

package storageaccounts

import (
	"azure-security-hub/utils"
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
)

// EnsureInfrastructureEncryption checks if infrastructure encryption is enabled for all storage accounts
func EnsureInfrastructureEncryption() string {
	// Load and print compliance info
	compliance, err := utils.LoadComplianceData("compliance/cis_microsoft_azure_foundations_benchmark_v3.0.0.json")
	if err != nil {
		log.Printf("└─[ERROR] Error loading compliance data: %v", err)
		return "NA"
	}
	utils.PrintComplianceInfo(compliance, "4.2")

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

			if account.Properties == nil ||
				account.Properties.Encryption == nil ||
				account.Properties.Encryption.RequireInfrastructureEncryption == nil ||
				!*account.Properties.Encryption.RequireInfrastructureEncryption {

				nonCompliantAccounts = append(nonCompliantAccounts, accountName)
				log.Printf("    └─[FAIL] Infrastructure encryption not enabled for: %s", accountName)

				// Additional encryption details
				if account.Properties != nil && account.Properties.Encryption != nil {
					log.Printf("      └─[INFO] Encryption settings:")
					if account.Properties.Encryption.KeySource != nil {
						log.Printf("        └─[INFO] Key Source: %s", *account.Properties.Encryption.KeySource)
					}
					if account.Properties.Encryption.RequireInfrastructureEncryption != nil {
						log.Printf("        └─[INFO] Infrastructure Encryption Required: %t",
							*account.Properties.Encryption.RequireInfrastructureEncryption)
					}
				}
			} else {
				log.Printf("    └─[PASS] Infrastructure encryption enabled for: %s", accountName)

				// Log encryption details for compliant accounts
				log.Printf("      └─[INFO] Encryption settings:")
				log.Printf("        └─[INFO] Key Source: %s", *account.Properties.Encryption.KeySource)
				log.Printf("        └─[INFO] Infrastructure Encryption Required: %t",
					*account.Properties.Encryption.RequireInfrastructureEncryption)
			}
		}
	}

	if accountsChecked == 0 {
		log.Println("└─[*] No storage accounts found")
		return "NA"
	}

	if len(nonCompliantAccounts) > 0 {
		log.Printf("└─[FAIL] Found %d storage accounts without infrastructure encryption:",
			len(nonCompliantAccounts))
		for _, account := range nonCompliantAccounts {
			log.Printf("  └─[FAIL] %s", account)
		}
		return "FAIL"
	}

	log.Printf("└─[PASS] All %d storage accounts have infrastructure encryption enabled", accountsChecked)
	return "PASS"
}
