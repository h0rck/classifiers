package ui

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"relatorios/services"
	"sort"
	"strings"
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
	fmt.Println("3. Mostrar regras de classificação atual")
	fmt.Println("4. Recarregar regras de classificação")
	fmt.Println("5. Sair")
	fmt.Println()

	fmt.Print("Digite sua escolha (1-5): ")
	choice, _ := ci.reader.ReadString('\n')
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
		fmt.Println("Encerrando o programa...")
		os.Exit(0)
	default:
		fmt.Println("\nOpção inválida! Pressione Enter para continuar...")
		ci.readLine()
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
	ci.readLine()
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
	ci.readLine()
	ci.showMainMenu()
}

// selectFile exibe uma navegação de arquivos para o usuário selecionar um arquivo
func (ci *ConsoleInterface) selectFile() {
	// Começar a partir do diretório atual
	startDir, err := os.Getwd()
	if err != nil {
		startDir = "/"
	}

	selectedPath, err := ci.browseFiles(startDir, true) // true significa que queremos selecionar um arquivo
	if err != nil {
		fmt.Printf("\nErro ao navegar pelos arquivos: %v\n", err)
		fmt.Print("\nPressione Enter para voltar ao menu principal...")
		ci.readLine()
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
	ci.readLine()
	ci.showMainMenu()
}

// selectDirectory exibe uma navegação de diretórios para o usuário selecionar uma pasta
func (ci *ConsoleInterface) selectDirectory() {
	// Começar a partir do diretório atual
	startDir, err := os.Getwd()
	if err != nil {
		startDir = "/"
	}

	selectedPath, err := ci.browseFiles(startDir, false) // false significa que queremos selecionar um diretório
	if err != nil {
		fmt.Printf("\nErro ao navegar pelos diretórios: %v\n", err)
		fmt.Print("\nPressione Enter para voltar ao menu principal...")
		ci.readLine()
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
	ci.readLine()
	ci.showMainMenu()
}

// browseFiles permite ao usuário navegar pelos arquivos e diretórios
func (ci *ConsoleInterface) browseFiles(currentDir string, selectFile bool) (string, error) {
	for {
		fmt.Print("\033[H\033[2J") // Limpa a tela

		// Mostrar diretório atual
		fmt.Printf("Diretório atual: %s\n\n", currentDir)

		if selectFile {
			fmt.Println("=== Selecione um ARQUIVO ou navegue pelos diretórios ===")
		} else {
			fmt.Println("=== Selecione um DIRETÓRIO ===")
		}

		fmt.Println("[0] .. (Voltar)")
		fmt.Println("[C] Cancelar")
		fmt.Println("[M] Digitar caminho manualmente")
		if !selectFile {
			fmt.Println("[S] Selecionar o diretório atual")
		}
		fmt.Println("-------------------------------------")

		// Ler os arquivos e diretórios
		entries, err := os.ReadDir(currentDir)
		if err != nil {
			return "", err
		}

		// Separar diretórios e arquivos
		var dirs []os.DirEntry
		var files []os.DirEntry

		// Filtrar arquivos compatíveis e diretórios
		for _, entry := range entries {
			if entry.IsDir() {
				dirs = append(dirs, entry)
			} else if selectFile && ci.isFormatSupported(entry.Name()) {
				// Se estamos selecionando um arquivo, mostrar apenas formatos suportados
				files = append(files, entry)
			}
		}

		// Ordenar os diretórios e arquivos alfabeticamente
		sort.Slice(dirs, func(i, j int) bool {
			return dirs[i].Name() < dirs[j].Name()
		})
		sort.Slice(files, func(i, j int) bool {
			return files[i].Name() < files[j].Name()
		})

		// Listar diretórios primeiro
		for i, entry := range dirs {
			fmt.Printf("[%d] 📁 %s\n", i+1, entry.Name())
		}

		// Então listar arquivos se estamos selecionando um arquivo
		if selectFile {
			for i, entry := range files {
				fmt.Printf("[%d] 📄 %s\n", i+len(dirs)+1, entry.Name())
			}
		}

		// Solicitar escolha
		fmt.Print("\nEscolha uma opção (número, C, M" + func() string {
			if !selectFile {
				return ", S"
			}
			return ""
		}() + "): ")

		input, _ := ci.readLine()
		input = strings.ToUpper(strings.TrimSpace(input))

		// Verificar opções especiais
		switch input {
		case "0": // Voltar
			parent := filepath.Dir(currentDir)
			if parent == currentDir { // Já estamos na raiz
				continue
			}
			currentDir = parent
			continue
		case "C": // Cancelar
			return "", nil
		case "M": // Digitar caminho manualmente
			fmt.Print("\nDigite o caminho completo: ")
			path, _ := ci.readLine()
			path = strings.TrimSpace(path)

			if path == "" {
				continue
			}

			// Verificar se o caminho existe
			info, err := os.Stat(path)
			if err != nil {
				fmt.Printf("\nErro: Caminho inválido ou inacessível.\n")
				fmt.Print("Pressione Enter para continuar...")
				ci.readLine()
				continue
			}

			// Se queremos um arquivo mas selecionamos um diretório, navegar para ele
			if selectFile && info.IsDir() {
				currentDir = path
				continue
			}

			// Se queremos um diretório e selecionamos um arquivo, mostrar erro
			if !selectFile && !info.IsDir() {
				fmt.Printf("\nErro: O caminho selecionado não é um diretório.\n")
				fmt.Print("Pressione Enter para continuar...")
				ci.readLine()
				continue
			}

			return path, nil
		case "S": // Selecionar diretório atual (apenas se !selectFile)
			if !selectFile {
				return currentDir, nil
			}
		}

		// Converter a escolha em número
		var index int
		if _, err := fmt.Sscanf(input, "%d", &index); err != nil {
			continue // Se não for um número, ignorar
		}

		// Verificar se é um diretório
		if index >= 1 && index <= len(dirs) {
			currentDir = filepath.Join(currentDir, dirs[index-1].Name())
			continue
		}

		// Verificar se é um arquivo e estamos selecionando um arquivo
		if selectFile && index > len(dirs) && index <= len(dirs)+len(files) {
			fileIndex := index - len(dirs) - 1
			selectedFile := filepath.Join(currentDir, files[fileIndex].Name())
			return selectedFile, nil
		}
	}
}

// isFormatSupported verifica se o formato do arquivo é suportado
func (ci *ConsoleInterface) isFormatSupported(filename string) bool {
	supportedFormats := ci.processingService.GetSupportedFormats()
	ext := strings.ToLower(filepath.Ext(filename))

	for _, format := range supportedFormats {
		if ext == format {
			return true
		}
	}
	return false
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
