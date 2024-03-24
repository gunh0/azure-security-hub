// audit/microsoftentraid/restrict_guest_user_access.go

package microsoftentraid

import (
	"azure-security-hub/audit"
	"azure-security-hub/utils"
	"context"
	"log"

	"github.com/google/uuid"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

const (
	// Role ID for most restrictive guest access
	restrictedGuestRoleID = "2af84b1e-32c8-42b7-82bc-daa82404023b"
)

// EnsureGuestUserAccessRestricted checks if guest user access is properly restricted
func EnsureGuestUserAccessRestricted() string {
	// Load and print compliance info
	compliance, err := utils.LoadComplianceData("compliance/cis_microsoft_azure_foundations_benchmark_v3.0.0.json")
	if err != nil {
		log.Printf("└─[ERROR] Error loading compliance data: %v", err)
		return "NA"
	}
	utils.PrintComplianceInfo(compliance, "2.15")

	// Get Azure client
	azureClient, err := utils.GetAzureClient()
	if err != nil {
		log.Printf("└─[ERROR] Failed to get Azure client: %v", err)
		return "NA"
	}

	log.Println("└─[*] Creating Microsoft Graph client")
	// Create Graph client with required scope
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

	guestRoleID := authPolicy.GetGuestUserRoleId()
	if guestRoleID == nil {
		log.Printf("└─[ERROR] Guest user role ID is not set")
		return "FAIL"
	}

	log.Printf("  └─[DEBUG] Current Guest User Role ID: %s", *guestRoleID)
	log.Printf("  └─[DEBUG] Expected Guest User Role ID: %s", restrictedGuestRoleID)

	restrictedGuestUUID, err := uuid.Parse(restrictedGuestRoleID)
	if err != nil {
		log.Printf("└─[ERROR] Failed to parse restricted guest role ID: %v", err)
		return "NA"
	}
	isRestricted := *guestRoleID == restrictedGuestUUID

	if isRestricted {
		log.Println("└─[PASS] Guest user access is properly restricted to their own directory objects")
		return "PASS"
	} else {
		log.Println("└─[FAIL] Guest user access is not set to most restrictive level")
		log.Println("└─[INFO] Guest access should be restricted to properties and memberships of their own directory objects")
		log.Printf("  └─[INFO] Current Guest User Role ID: %s", *guestRoleID)
		log.Printf("  └─[INFO] Expected Guest User Role ID: %s", restrictedGuestRoleID)
		return "FAIL"
	}
}
