// audit/common.go
package audit

import (
	"log"

	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// Helper function to convert boolean condition to status string
func GetStatus(condition bool) string {
	if condition {
		return "PASS"
	}
	return "FAIL"
}

// Helper function to print detailed OData errors
func PrintOdataError(err error) {
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
