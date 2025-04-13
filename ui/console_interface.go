package ui

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"relatorios/models"
	"relatorios/services"
	"strings"
)

// ConsoleInterface gerencia a interface de linha de comando da aplicação
type ConsoleInterface struct {
	processingService *services.DocumentProcessingService
	reader            *bufio.Reader
	fileBrowser       *FileBrowser
}

// NewConsoleInterface cria uma nova instância da interface de console
func NewConsoleInterface(processingService *services.DocumentProcessingService) *ConsoleInterface {
	consoleInterface := &ConsoleInterface{
		processingService: processingService,
		reader:            bufio.NewReader(os.Stdin),
	}

	// Create file browser with console interface as the delegate
	consoleInterface.fileBrowser = NewFileBrowser(consoleInterface, consoleInterface)

	return consoleInterface
}

// ReadLine implements InputReader interface
func (ci *ConsoleInterface) ReadLine() (string, error) {
	text, err := ci.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}

// IsFormatSupported implements FormatChecker interface
func (ci *ConsoleInterface) IsFormatSupported(filename string) bool {
	supportedFormats := ci.processingService.GetSupportedFormats()
	ext := strings.ToLower(filepath.Ext(filename))

	for _, format := range supportedFormats {
		if ext == format {
			return true
		}
	}
	return false
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

	// Mostrar o arquivo de regras atual
	analyzeService := ci.processingService.GetAnalyzeService()
	fmt.Printf("Arquivo de regras: %s\n", analyzeService.GetRulesFilePath())

	fmt.Println("\nSelecione uma opção:")
	fmt.Println("1. Processar um único arquivo")
	fmt.Println("2. Processar todos os arquivos em uma pasta")
	fmt.Println("3. Mostrar regras de classificação atual")
	fmt.Println("4. Recarregar regras de classificação")
	fmt.Println("5. Selecionar arquivo de regras")
	fmt.Println("6. Sair")
	fmt.Println()

	fmt.Print("Digite sua escolha (1-6): ")
	choice, _ := ci.ReadLine()
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		ci.selectFile()
	case "2":
		ci.selectDirectory()
	case "3":
		ci.showClassificationRules()
	case "4":
		ci.reloadClassificationRules()
	case "5":
		ci.selectRulesFile()
	case "6":
		fmt.Println("Encerrando o programa...")
		os.Exit(0)
	default:
		fmt.Println("\nOpção inválida! Pressione Enter para continuar...")
		ci.ReadLine()
		ci.showMainMenu()
	}
}

// showClassificationRules exibe as regras de classificação atuais
func (ci *ConsoleInterface) showClassificationRules() {
	fmt.Print("\033[H\033[2J") // Limpa a tela
	fmt.Println("=== Regras de Classificação Atuais ===")

	// Obter as regras do serviço
	analyzeService := ci.processingService.GetAnalyzeService()
	rules := analyzeService.GetRules()

	// Mostrar o caminho do arquivo de regras
	fmt.Printf("Arquivo de regras: %s\n\n", analyzeService.GetRulesFilePath())

	// Exibir cada regra
	for i, rule := range rules {
		fmt.Printf("%d. Tipo: %s\n", i+1, rule.Type)
		fmt.Printf("   Palavras-chave: %s\n\n", strings.Join(rule.Keywords, ", "))
	}

	fmt.Print("\nPressione Enter para voltar ao menu principal...")
	ci.ReadLine()
	ci.showMainMenu()
}

// reloadClassificationRules recarrega as regras de classificação do arquivo JSON
func (ci *ConsoleInterface) reloadClassificationRules() {
	fmt.Print("\033[H\033[2J") // Limpa a tela
	fmt.Println("=== Recarregando Regras de Classificação ===")

	// Recarregar as regras
	analyzeService := ci.processingService.GetAnalyzeService()
	err := analyzeService.ReloadRules()

	if err != nil {
		fmt.Printf("\nErro ao recarregar regras: %v\n", err)
	} else {
		fmt.Printf("\nRegras recarregadas com sucesso do arquivo:\n%s\n", analyzeService.GetRulesFilePath())
	}

	fmt.Print("\nPressione Enter para voltar ao menu principal...")
	ci.ReadLine()
	ci.showMainMenu()
}

// selectRulesFile permite ao usuário selecionar um novo arquivo de regras JSON
func (ci *ConsoleInterface) selectRulesFile() {
	// Começar a partir do diretório atual
	startDir, err := os.Getwd()
	if err != nil {
		startDir = "/"
	}

	fmt.Print("\033[H\033[2J") // Limpa a tela
	fmt.Println("=== Selecionar Arquivo de Regras JSON ===")
	fmt.Println("Navegue até o arquivo JSON que contém suas regras de classificação.")
	fmt.Print("\nPressione Enter para continuar...")
	ci.ReadLine()

	// Usar o navegador de arquivos com filtro para arquivos JSON
	selectedPath, err := ci.fileBrowser.BrowseFilesWithFilter(startDir, true, []string{".json"})
	if err != nil {
		fmt.Printf("\nErro ao navegar pelos arquivos: %v\n", err)
		fmt.Print("\nPressione Enter para voltar ao menu principal...")
		ci.ReadLine()
		ci.showMainMenu()
		return
	}

	// Se o usuário cancelou a seleção
	if selectedPath == "" {
		ci.showMainMenu()
		return
	}

	// Atualizar o arquivo de regras
	analyzeService := ci.processingService.GetAnalyzeService()

	// Tentar carregar as regras do novo arquivo
	newRules, err := models.LoadRulesFromJSON(selectedPath)
	if err != nil {
		fmt.Printf("\nErro ao carregar regras do arquivo: %v\n", err)
		fmt.Print("\nPressione Enter para voltar ao menu principal...")
		ci.ReadLine()
		ci.showMainMenu()
		return
	}

	// Se foi bem sucedido, atualizar o serviço
	err = analyzeService.SetRulesFile(selectedPath)
	if err != nil {
		fmt.Printf("\nErro ao definir novo arquivo de regras: %v\n", err)
	} else {
		fmt.Printf("\nArquivo de regras atualizado para: %s\n", selectedPath)
		fmt.Printf("Carregadas %d regras de classificação\n", len(newRules))
	}

	fmt.Print("\nPressione Enter para voltar ao menu principal...")
	ci.ReadLine()
	ci.showMainMenu()
}

// selectFile exibe uma navegação de arquivos para o usuário selecionar um arquivo
func (ci *ConsoleInterface) selectFile() {
	// Começar a partir do diretório atual
	startDir, err := os.Getwd()
	if err != nil {
		startDir = "/"
	}

	selectedPath, err := ci.fileBrowser.BrowseFiles(startDir, true) // true significa que queremos selecionar um arquivo
	if err != nil {
		fmt.Printf("\nErro ao navegar pelos arquivos: %v\n", err)
		fmt.Print("\nPressione Enter para voltar ao menu principal...")
		ci.ReadLine()
		ci.showMainMenu()
		return
	}

	// Se o usuário cancelou a seleção
	if selectedPath == "" {
		ci.showMainMenu()
		return
	}

	// Processar o arquivo selecionado
	if err := ci.handleSingleFile(selectedPath); err != nil {
		fmt.Printf("\nErro: %v\n", err)
	}

	fmt.Print("\nPressione Enter para voltar ao menu principal...")
	ci.ReadLine()
	ci.showMainMenu()
}

// selectDirectory exibe uma navegação de diretórios para o usuário selecionar uma pasta
func (ci *ConsoleInterface) selectDirectory() {
	// Começar a partir do diretório atual
	startDir, err := os.Getwd()
	if err != nil {
		startDir = "/"
	}

	selectedPath, err := ci.fileBrowser.BrowseFiles(startDir, false) // false significa que queremos selecionar um diretório
	if err != nil {
		fmt.Printf("\nErro ao navegar pelos diretórios: %v\n", err)
		fmt.Print("\nPressione Enter para voltar ao menu principal...")
		ci.ReadLine()
		ci.showMainMenu()
		return
	}

	// Se o usuário cancelou a seleção
	if selectedPath == "" {
		ci.showMainMenu()
		return
	}

	// Processar o diretório selecionado
	if err := ci.handleDirectory(selectedPath); err != nil {
		fmt.Printf("\nErro: %v\n", err)
	}

	fmt.Print("\nPressione Enter para voltar ao menu principal...")
	ci.ReadLine()
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
