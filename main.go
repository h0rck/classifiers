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
	// Definir o caminho para o arquivo de regras
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		userConfigDir = "."
	}

	configDir := filepath.Join(userConfigDir, "relatorios-go")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Aviso: Não foi possível criar diretório de configuração: %v\n", err)
		configDir = "."
	}

	rulesFile := filepath.Join(configDir, "document_rules.json")

	// Inicializar componentes do sistema
	extractorFactory := extractors.NewDocumentExtractorFactory()
	analyzeDocumentService := services.NewAnalyzeDocumentService(rulesFile)
	classifier := classifiers.NewDocumentClassifier(analyzeDocumentService)

	// Configurar o serviço de processamento de documentos
	config := models.ProcessingConfig{
		OutputDirectory: "./output",
		MoveFiles:       false,
	}

	processingService := services.NewDocumentProcessingService(
		extractorFactory,
		classifier,
		config,
	)

	// Criar a interface de console
	consoleInterface := ui.NewConsoleInterface(processingService)

	// Determinar se o programa foi chamado com um argumento de caminho
	var initialPath string
	if len(os.Args) > 1 {
		initialPath = os.Args[1]
	}

	// Iniciar a interface
	err = consoleInterface.Start(initialPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao iniciar a aplicação: %v\n", err)
		os.Exit(1)
	}
}
