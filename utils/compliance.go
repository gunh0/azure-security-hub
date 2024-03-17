package utils

import (
	"encoding/json"
	"log"
	"os"
)

type Attribute struct {
	Section1           string `json:"Section1"`
	Section2           string `json:"Section2"`
	AssessmentStatus   string `json:"AssessmentStatus"`
	ApplicableProfiles string `json:"ApplicableProfiles"`
	Description        string `json:"Description"`
	RationaleStatement string `json:"RationaleStatement"`
	ImpactStatement    string `json:"ImpactStatement"`
}

type Requirement struct {
	Id         string      `json:"Id"`
	Title      string      `json:"Title"`
	Checks     []string    `json:"Checks"`
	Attributes []Attribute `json:"Attributes"`
}

type Compliance struct {
	Framework    string        `json:"Framework"`
	Version      string        `json:"Version"`
	Provider     string        `json:"Provider"`
	Description  string        `json:"Description"`
	Requirements []Requirement `json:"Requirements"`
}

func LoadComplianceData(filePath string) (*Compliance, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var compliance Compliance
	if err := json.Unmarshal(data, &compliance); err != nil {
		return nil, err
	}

	return &compliance, nil
}

// Print the specific compliance information
func PrintComplianceInfo(compliance *Compliance, id string) {
	for _, requirement := range compliance.Requirements {
		if requirement.Id == id {
			log.Printf("[%s %s %s] %s - %s\n",
				compliance.Framework,
				compliance.Provider,
				compliance.Version,
				requirement.Id,
				requirement.Title)
			return
		}
	}
	log.Printf("[ERROR] Compliance requirement with ID %s not found.\n", id)
}
