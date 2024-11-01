package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"azure-security-hub/audit/microsoftentraid"
	"azure-security-hub/audit/storageaccounts"
	"azure-security-hub/utils"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "azure-security-hub",
	Short: "Azure Security Hub CLI",
	Long:  `A CLI tool for Azure Security Hub operations.`,
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Azure",
	Long:  `Authenticate with Azure using the credentials in your .env file.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[*] Attempting to login to Azure...")
		client, err := utils.LoginAzure()
		if err != nil {
			fmt.Printf("[-] Failed to login: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("[+] Login successful!")

		// Use the client to list subscriptions
		fmt.Println("[*] Listing subscriptions:")
		pager := client.NewListPager(nil)
		for pager.More() {
			page, err := pager.NextPage(context.Background())
			if err != nil {
				fmt.Printf("[-] Failed to list subscriptions: %v\n", err)
				os.Exit(1)
			}
			for _, sub := range page.Value {
				fmt.Printf("[+] %s (%s)\n", *sub.DisplayName, *sub.SubscriptionID)
			}
		}
		fmt.Println("[*] Subscription listing complete.")
	},
}

// 2.3
var restrictTenantCreationCmd = &cobra.Command{
	Use:     "restrict-tenant-creation",
	Short:   "Ensure that 'Restrict non-admin users from creating tenants' is set to 'Yes'",
	Aliases: []string{"cis.3.0.0-identity.2.3"},
	Run: func(cmd *cobra.Command, args []string) {
		result := microsoftentraid.EnsureTenantCreationRestricted()
		log.Printf("[Microsoft Entra ID] (%s) : %s", result, cmd.Short)
	},
}

// 2.14
var restrictAppRegistrationCmd = &cobra.Command{
	Use:     "restrict-app-registration",
	Short:   "Ensure That 'Users Can Register Application' Is Set to 'No'",
	Aliases: []string{"cis.3.0.0-identity.2.14"},
	Run: func(cmd *cobra.Command, args []string) {
		result := microsoftentraid.EnsureAppRegistrationRestricted()
		log.Printf("[Microsoft Entra ID] (%s) : %s", result, cmd.Short)
	},
}

// 2.15
var restrictGuestUserAccessCmd = &cobra.Command{
	Use:     "restrict-guest-user-access",
	Short:   "Ensure That 'Guest users access restrictions' is set to 'Guest user access is restricted to properties and memberships of their own directory objects'",
	Aliases: []string{"cis.3.0.0-identity.2.15"},
	Run: func(cmd *cobra.Command, args []string) {
		result := microsoftentraid.EnsureGuestUserAccessRestricted()
		log.Printf("[Microsoft Entra ID] (%s) : %s", result, cmd.Short)
	},
}

// 4.1
var ensureSecureTransferRequiredCmd = &cobra.Command{
	Use:     "ensure-secure-transfer-required",
	Short:   "Ensure that 'Secure transfer required' is set to 'Enabled'",
	Aliases: []string{"cis.3.0.0-storage-accounts.4.1"},
	Run: func(cmd *cobra.Command, args []string) {
		result := storageaccounts.EnsureSecureTransferRequired()
		log.Printf("[Storage Accounts] (%s) : %s", result, cmd.Short)
	},
}

// 4.2
var ensureInfrastructureEncryptionCmd = &cobra.Command{
	Use:     "ensure-infrastructure-encryption",
	Short:   "Ensure that 'Enable Infrastructure Encryption' for Storage Accounts is Set to 'enabled'",
	Aliases: []string{"storage.4.2"},
	Run: func(cmd *cobra.Command, args []string) {
		result := storageaccounts.EnsureInfrastructureEncryption()
		log.Printf("[Storage Accounts] (%s) : %s", result, cmd.Short)
	},
}

// init is called before the main function
func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(restrictTenantCreationCmd)         // 2.3
	rootCmd.AddCommand(restrictAppRegistrationCmd)        // 2.14
	rootCmd.AddCommand(restrictGuestUserAccessCmd)        // 2.15
	rootCmd.AddCommand(ensureSecureTransferRequiredCmd)   // 4.1
	rootCmd.AddCommand(ensureInfrastructureEncryptionCmd) // 4.2
}

// main is the entry point of the application
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("[-] Error: %v\n", err)
		os.Exit(1)
	}
}
