package main

import (
	"fmt"
	"os"
	"relatorios/models"
	"relatorios/services"
	"relatorios/services/classifiers"
	"relatorios/services/extractors"
	"relatorios/ui"
)

func main() {
	// Inicializar componentes do sistema
	extractorFactory := extractors.NewDocumentExtractorFactory()
	analyzeDocumentService := services.NewAnalyzeDocumentService()
	classifier := classifiers.NewDocumentClassifier(analyzeDocumentService)

	// Configurar o serviço de processamento de documentos
	config := models.ProcessingConfig{
		OutputDirectory: "./output",
		// aqui ele cria uma copia do arquivo
		MoveFiles: false,
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
	err := consoleInterface.Start(initialPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao iniciar a aplicação: %v\n", err)
		os.Exit(1)
	}
}
