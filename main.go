package main

import (
	"fmt"
	"os"
	"path/filepath"
	"relatorios/models"
	"relatorios/services"
	"relatorios/services/classifiers"
	"relatorios/services/extractors"
	"relatorios/ui"
)

func main() {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		userConfigDir = "."
	}

	configDir := filepath.Join(userConfigDir, "relatorios-go")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not create configuration directory: %v\n", err)
		configDir = "."
	}

	rulesFile := filepath.Join(configDir, "document_rules.json")

	extractorFactory := extractors.NewDocumentExtractorFactory()
	analyzeDocumentService := services.NewAnalyzeDocumentService(rulesFile)
	classifier := classifiers.NewDocumentClassifier(analyzeDocumentService)

	config := models.ProcessingConfig{
		OutputDirectory: "./output",
		MoveFiles:       false,
	}

	processingService := services.NewDocumentProcessingService(
		extractorFactory,
		classifier,
		config,
	)

	consoleInterface := ui.NewConsoleInterface(processingService)

	var initialPath string
	if len(os.Args) > 1 {
		initialPath = os.Args[1]
	}

	err = consoleInterface.Start(initialPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting application: %v\n", err)
		os.Exit(1)
	}
}
