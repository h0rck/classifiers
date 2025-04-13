package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"relatorios/services"
)

// ConsoleInterface gerencia a interface de linha de comando da aplicação
type ConsoleInterface struct {
	processingService *services.DocumentProcessingService
	reader            *bufio.Reader
}

// NewConsoleInterface cria uma nova instância da interface de console
func NewConsoleInterface(processingService *services.DocumentProcessingService) *ConsoleInterface {
	return &ConsoleInterface{
		processingService: processingService,
		reader:            bufio.NewReader(os.Stdin),
	}
}

// Start inicia a interface do console com um caminho opcional
func (ci *ConsoleInterface) Start(initialPath string) error {
	if initialPath != "" {
		return ci.handleInitialPath(initialPath)
	}

	ci.showMainMenu()
	return nil
}

// handleInitialPath processa o caminho inicial fornecido pelo usuário
func (ci *ConsoleInterface) handleInitialPath(pathInput string) error {
	fileInfo, err := os.Stat(pathInput)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("\nCaminho não encontrado: %s\n", pathInput)
			ci.showMainMenu()
			return nil
		}
		return fmt.Errorf("erro ao verificar caminho: %w", err)
	}

	if fileInfo.IsDir() {
		return ci.handleDirectory(pathInput)
	}

	return ci.handleSingleFile(pathInput)
}

// showMainMenu exibe o menu principal da aplicação
func (ci *ConsoleInterface) showMainMenu() {
	fmt.Print("\033[H\033[2J") // Limpa a tela (equivalente ao console.clear())
	fmt.Println("====================================")
	fmt.Println("=== CLASSIFICADOR DE DOCUMENTOS ===")
	fmt.Println("====================================")
	fmt.Printf("\nClassificador: %s\n", ci.processingService.GetClassifierName())
	fmt.Printf("Formatos suportados: %s\n", strings.Join(ci.processingService.GetSupportedFormats(), ", "))
	fmt.Println("\nSelecione uma opção:")
	fmt.Println("1. Processar um único arquivo")
	fmt.Println("2. Processar todos os arquivos em uma pasta")
	fmt.Println("3. Sair")
	fmt.Println()

	fmt.Print("Digite sua escolha (1-3): ")
	choice, _ := ci.reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		ci.askForFilePath()
	case "2":
		ci.askForDirectoryPath()
	case "3":
		fmt.Println("Encerrando o programa...")
		os.Exit(0)
	default:
		fmt.Println("\nOpção inválida! Pressione Enter para continuar...")
		ci.readLine()
		ci.showMainMenu()
	}
}

// askForFilePath solicita ao usuário o caminho de um arquivo
func (ci *ConsoleInterface) askForFilePath() {
	fmt.Print("\nDigite o caminho completo do arquivo a ser classificado: ")
	filePath, _ := ci.readLine()

	if err := ci.handleSingleFile(filePath); err != nil {
		fmt.Printf("\nErro: %v\n", err)
	}

	fmt.Print("\nPressione Enter para voltar ao menu principal...")
	ci.readLine()
	ci.showMainMenu()
}

// askForDirectoryPath solicita ao usuário o caminho de um diretório
func (ci *ConsoleInterface) askForDirectoryPath() {
	fmt.Print("\nDigite o caminho completo da pasta com os documentos: ")
	dirPath, _ := ci.readLine()

	if err := ci.handleDirectory(dirPath); err != nil {
		fmt.Printf("\nErro: %v\n", err)
	}

	fmt.Print("\nPressione Enter para voltar ao menu principal...")
	ci.readLine()
	ci.showMainMenu()
}

// handleSingleFile processa um único arquivo
func (ci *ConsoleInterface) handleSingleFile(filePath string) error {
	fmt.Printf("\nProcessando o arquivo: %s\n", filePath)

	document, destinationPath, err := ci.processingService.ProcessSingleFile(filePath)
	if err != nil {
		return err
	}

	fmt.Println("\n===== Resultado da Classificação =====")
	fmt.Printf("Arquivo: %s\n", document.Filename)
	if document.Classification != nil {
		fmt.Printf("Tipo de documento: %s\n", document.Classification.DocumentType)
		fmt.Printf("Palavras-chave: %s\n", strings.Join(document.Classification.Keywords, ", "))
	} else {
		fmt.Println("Não foi possível classificar o documento")
	}
	fmt.Printf("\nArquivo organizado em: %s\n", destinationPath)

	return nil
}

// handleDirectory processa todos os arquivos em um diretório
func (ci *ConsoleInterface) handleDirectory(dirPath string) error {
	fmt.Printf("\nProcessando diretório: %s\n", dirPath)

	result, err := ci.processingService.ProcessDirectory(dirPath)
	if err != nil {
		return err
	}

	fmt.Println("\n===== Resultado do Processamento =====")
	fmt.Printf("Total de arquivos processados: %d\n", result.ProcessedCount)
	fmt.Printf("Total de falhas: %d\n", result.FailedCount)

	for _, fileResult := range result.Results {
		if fileResult.Success {
			fmt.Printf("\n✓ %s → %s\n",
				fileResult.Filename,
				fileResult.DocumentType)
		} else {
			fmt.Printf("\n✗ %s → FALHA: %s\n",
				fileResult.Filename,
				fileResult.Error)
		}
	}

	fmt.Printf("\nDocumentos organizados em: %s\n", ci.processingService.GetOutputDirectory())

	return nil
}

// readLine lê uma linha do entrada padrão e remove os espaços em branco
func (ci *ConsoleInterface) readLine() (string, error) {
	text, err := ci.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}
