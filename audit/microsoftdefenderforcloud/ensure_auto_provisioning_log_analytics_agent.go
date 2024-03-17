package microsoftdefenderforcloud

import (
	"azure-security-hub/utils"
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/security/armsecurity"
)

var newAutoProvisioningSettingsClient = armsecurity.NewAutoProvisioningSettingsClient

func EnsureAutoProvisioningLogAnalyticsAgent() string {
	// Load and print compliance info
	compliance, err := utils.LoadComplianceData("compliance/cis_microsoft_azure_foundations_benchmark_v3.0.0.json")
	if err != nil {
		log.Printf("└─[ERROR] Error loading compliance data: %v", err)
		return "NA"
	}
	utils.PrintComplianceInfo(compliance, "3.1.1.1")

	// Get Azure client
	client, err := utils.GetAzureClient()
	if err != nil {
		log.Printf("└─[ERROR] Failed to get Azure client: %v", err)
		return "NA"
	}

	// Create auto provisioning settings client
	log.Println("└─[*] Creating auto provisioning settings client")
	autoProvisioningClient, err := newAutoProvisioningSettingsClient(client.SubscriptionID, client.Cred, nil)
	if err != nil {
		log.Printf("└─[ERROR] Failed to create auto provisioning settings client: %v", err)
		return "NA"
	}

	// Get auto provisioning settings
	log.Println("└─[*] Checking auto provisioning settings")
	defaultSettingFound := false
	pager := autoProvisioningClient.NewListPager(nil)
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			log.Printf("└─[ERROR] Failed to get auto provisioning settings: %v", err)
			return "NA"
		}

		for _, setting := range page.Value {
			if setting.Name == nil {
				continue
			}

			log.Printf("  └─[*] Checking setting: %s", *setting.Name)
			if *setting.Name == "default" {
				defaultSettingFound = true
				if setting.Properties == nil || setting.Properties.AutoProvision == nil {
					log.Println("    └─[FAIL] Auto provisioning property is not set")
					return "FAIL"
				}

				autoProvision := *setting.Properties.AutoProvision
				log.Printf("    └─[%s] Auto provisioning: %s", getStatus(autoProvision == armsecurity.AutoProvisionOn), string(autoProvision))

				if autoProvision == armsecurity.AutoProvisionOn {
					log.Println("└─[PASS] Auto provisioning of 'Log Analytics agent for Azure VMs' is set to 'On'")
					return "PASS"
				} else {
					log.Println("└─[FAIL] Auto provisioning of 'Log Analytics agent for Azure VMs' is not set to 'On'")
					return "FAIL"
				}
			}
		}
	}

	if !defaultSettingFound {
		log.Println("└─[FAIL] Default auto provisioning setting not found")
	}
	return "FAIL"
}

func getStatus(condition bool) string {
	if condition {
		return "PASS"
	}
	return "FAIL"
}
