// audit/microsoftentraid/restrict_app_registration.go

package microsoftentraid

import (
	"azure-security-hub/audit"
	"azure-security-hub/utils"
	"context"
	"log"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

func EnsureAppRegistrationRestricted() string {
	// Load and print compliance info
	compliance, err := utils.LoadComplianceData("compliance/cis_microsoft_azure_foundations_benchmark_v3.0.0.json")
	if err != nil {
		log.Printf("└─[ERROR] Error loading compliance data: %v", err)
		return "NA"
	}
	utils.PrintComplianceInfo(compliance, "2.14")

	// Get Azure client
	azureClient, err := utils.GetAzureClient()
	if err != nil {
		log.Printf("└─[ERROR] Failed to get Azure client: %v", err)
		return "NA"
	}

	log.Println("└─[*] Creating Microsoft Graph client")
	// Create Graph client with required permissions
	// Use .default scope for client credentials flow
	scopes := []string{"https://graph.microsoft.com/.default"}
	client, err := msgraphsdk.NewGraphServiceClientWithCredentials(azureClient.Cred, scopes)
	if err != nil {
		log.Printf("└─[ERROR] Error creating Graph client: %v", err)
		return "NA"
	}

	log.Println("└─[*] Requesting authorization policy")
	authPolicy, err := client.Policies().AuthorizationPolicy().Get(context.Background(), nil)
	if err != nil {
		log.Printf("└─[ERROR] Failed to get authorization policy: %v", err)
		audit.PrintOdataError(err)
		return "NA"
	}

	if authPolicy == nil {
		log.Printf("└─[ERROR] Authorization policy is nil")
		return "NA"
	}

	log.Printf("└─[DEBUG] Authorization Policy Response:")
	log.Printf("  └─[DEBUG] ID: %s", *authPolicy.GetId())
	log.Printf("  └─[DEBUG] DisplayName: %s", *authPolicy.GetDisplayName())

	defaultPermissions := authPolicy.GetDefaultUserRolePermissions()
	if defaultPermissions == nil {
		log.Printf("└─[ERROR] Default user role permissions is nil")
		return "NA"
	}

	allowedToCreateApps := defaultPermissions.GetAllowedToCreateApps()
	log.Printf("  └─[DEBUG] AllowedToCreateApps: %v", allowedToCreateApps)

	isRestricted := !*allowedToCreateApps

	log.Printf("  └─[*] Checking application registration restriction")
	log.Printf("    └─[%s] Users can register applications: %t",
		audit.GetStatus(!isRestricted), *allowedToCreateApps)

	if isRestricted {
		log.Println("└─[PASS] Users cannot register applications")
		return "PASS"
	} else {
		log.Println("└─[FAIL] Users are allowed to register applications")
		return "FAIL"
	}
}
