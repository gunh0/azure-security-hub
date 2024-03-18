package microsoftentraid

import (
	"azure-security-hub/utils"
	"context"
	"log"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

func EnsureTenantCreationRestricted() string {
	// Load and print compliance info
	compliance, err := utils.LoadComplianceData("compliance/cis_microsoft_azure_foundations_benchmark_v3.0.0.json")
	if err != nil {
		log.Printf("└─[ERROR] Error loading compliance data: %v", err)
		return "NA"
	}
	utils.PrintComplianceInfo(compliance, "2.3")

	// Get Azure client
	azureClient, err := utils.GetAzureClient()
	if err != nil {
		log.Printf("└─[ERROR] Failed to get Azure client: %v", err)
		return "NA"
	}

	log.Println("└─[*] Creating Microsoft Graph client")
	// Create Graph client with credentials and required scope
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
		printOdataError(err) // Print detailed error information
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

	allowedToCreateTenants := defaultPermissions.GetAllowedToCreateTenants()
	log.Printf("  └─[DEBUG] AllowedToCreateTenants: %v", allowedToCreateTenants)

	isRestricted := !*allowedToCreateTenants

	log.Printf("  └─[*] Checking tenant creation restriction")
	log.Printf("    └─[%s] Restrict non-admin users from creating tenants: %t",
		getStatus(isRestricted), isRestricted)

	if isRestricted {
		log.Println("└─[PASS] Non-admin users are restricted from creating tenants")
		return "PASS"
	} else {
		log.Println("└─[FAIL] Non-admin users are allowed to create tenants")
		return "FAIL"
	}
}

// Helper function to print detailed OData errors
func printOdataError(err error) {
	if odataErr, ok := err.(*odataerrors.ODataError); ok {
		log.Printf("└─[ERROR] OData Error: %v", odataErr)

		if odataErr.Error() != "" {
			log.Printf("  └─[ERROR] Details: %s", odataErr.Error())
		}

		// Additional message for permission errors
		if odataErr.ResponseStatusCode == 403 {
			log.Printf("└─[ERROR] Permission denied: 'Policy.Read.All' permission is required.")
		}

		// Additional error information if available in raw format
		log.Printf("  └─[ERROR] Raw Error Info: %#v", odataErr)
	} else {
		log.Printf("└─[ERROR] Non-OData error: %v", err)
	}
}

func getStatus(condition bool) string {
	if condition {
		return "PASS"
	}
	return "FAIL"
}
