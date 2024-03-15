package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
)

func TestLoginAzure(t *testing.T) {
	t.Log("[*] Starting LoginAzure test")

	// Load .env file for testing
	envPath := filepath.Join("..", ".env")
	err := godotenv.Load(envPath)
	if err != nil {
		t.Fatalf("[-] Error loading .env file: %v", err)
	}
	t.Log("[+] .env file loaded successfully")

	// Check if required environment variables are set
	requiredEnvVars := []string{"AZURE_TENANT_ID", "AZURE_CLIENT_ID", "AZURE_CLIENT_SECRET"}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			t.Fatalf("[-] %s is not set in the environment", envVar)
		}
	}
	t.Log("[+] All required environment variables are set")

	// Attempt to login
	client, err := LoginAzure()
	if err != nil {
		t.Fatalf("[-] LoginAzure failed: %v", err)
	}

	// Check if client is not nil
	if client == nil {
		t.Fatal("[-] LoginAzure returned nil client")
	}

	t.Log("[+] LoginAzure succeeded")
}

func TestLoginAzureWithInvalidCredentials(t *testing.T) {
	t.Log("[*] Starting LoginAzure with invalid credentials test")

	// Backup original environment variables
	originalTenantID := os.Getenv("AZURE_TENANT_ID")
	originalClientID := os.Getenv("AZURE_CLIENT_ID")
	originalClientSecret := os.Getenv("AZURE_CLIENT_SECRET")

	// Restore original environment variables after the test
	defer func() {
		os.Setenv("AZURE_TENANT_ID", originalTenantID)
		os.Setenv("AZURE_CLIENT_ID", originalClientID)
		os.Setenv("AZURE_CLIENT_SECRET", originalClientSecret)
		t.Log("[+] Original environment variables restored")
	}()

	// Set invalid credentials
	os.Setenv("AZURE_TENANT_ID", "invalid-tenant-id")
	os.Setenv("AZURE_CLIENT_ID", "invalid-client-id")
	os.Setenv("AZURE_CLIENT_SECRET", "invalid-client-secret")
	t.Log("[+] Invalid credentials set")

	// Attempt to login with invalid credentials
	_, err := LoginAzure()
	if err == nil {
		t.Fatal("[-] LoginAzure succeeded with invalid credentials")
	}

	t.Log("[+] LoginAzure correctly failed with invalid credentials")
}
