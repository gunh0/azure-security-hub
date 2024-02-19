package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
)

func main() {
	// Check if Azure CLI is installed
	if _, err := exec.LookPath("az"); err != nil {
		log.Fatal("[-] Azure CLI is not installed. Please install it first.")
	}

	// Delete existing Azure CLI configuration
	fmt.Println("[*] Deleting existing Azure CLI configuration...")
	exec.Command("make", "logout").Run()

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("[-] .env file does not exist. Skipping environment variables setup.")
	}

	// Check .env file for required variables
	clientID := os.Getenv("AZURE_CLIENT_ID")
	tenantID := os.Getenv("AZURE_TENANT_ID")
	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")

	if clientID == "" || tenantID == "" || clientSecret == "" {
		log.Fatal("[-] Required environment variables are not set. Please set them in the .env file.")
	}

	// Login to Azure CLI
	fmt.Println("[*] Using client secret for Azure CLI login.")
	fmt.Printf("  [*] AZURE_CLIENT_ID: %s\n", clientID)
	fmt.Printf("  [*] AZURE_TENANT_ID: %s\n", tenantID)

	cmd := exec.Command("az", "login", "--service-principal", "-u", clientID, "-p", clientSecret, "--tenant", tenantID)

	if output, err := cmd.CombinedOutput(); err != nil {
		log.Fatalf("[-] Azure CLI login failed: %v\n%s", err, output)
	} else {
		fmt.Println("[+] Azure CLI login success.")
	}
}
