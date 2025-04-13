package models

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type DocumentRule struct {
	Type     string   `json:"type"`
	Keywords []string `json:"keywords"`
}

func LoadRulesFromJSON(filePath string) ([]DocumentRule, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("\n===============================================================\n")
		fmt.Printf("WARNING: Rules file not found: %s\n\n", filePath)
		fmt.Printf("To use the classifier, create a JSON file with the following format:\n\n")
		fmt.Printf("[\n")
		fmt.Printf("  {\n")
		fmt.Printf("    \"type\": \"Document Type Name\",\n")
		fmt.Printf("    \"keywords\": [\n")
		fmt.Printf("      \"keyword1\",\n")
		fmt.Printf("      \"keyword2\",\n")
		fmt.Printf("      \"key phrase also works\"\n")
		fmt.Printf("    ]\n")
		fmt.Printf("  },\n")
		fmt.Printf("]\n\n")

		simpleRules := []DocumentRule{
			{
				Type:     "Document",
				Keywords: []string{"text", "document"},
			},
		}

		fmt.Printf("Creating a basic file to get you started...\n")

		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return simpleRules, nil
		}

		if err := SaveRulesToJSON(filePath, simpleRules); err != nil {
			fmt.Printf("Could not create rules file. Using minimal rules.\n")
			return simpleRules, nil
		}

		fmt.Printf("Basic rules file created at: %s\n", filePath)
		fmt.Printf("Edit this file to customize your classification.\n")
		fmt.Printf("===============================================================\n")

		return simpleRules, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read rules file: %w", err)
	}

	var rules []DocumentRule
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, fmt.Errorf("failed to decode JSON rules: %w", err)
	}

	return rules, nil
}

func SaveRulesToJSON(filePath string, rules []DocumentRule) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory for rules file: %w", err)
	}

	data, err := json.MarshalIndent(rules, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode rules to JSON: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to save rules file: %w", err)
	}

	return nil
}
