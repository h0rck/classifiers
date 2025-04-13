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

type ConsoleInterface struct {
	processingService *services.DocumentProcessingService
	reader            *bufio.Reader
	fileBrowser       *FileBrowser
}

func NewConsoleInterface(processingService *services.DocumentProcessingService) *ConsoleInterface {
	consoleInterface := &ConsoleInterface{
		processingService: processingService,
		reader:            bufio.NewReader(os.Stdin),
	}

	consoleInterface.fileBrowser = NewFileBrowser(consoleInterface, consoleInterface)

	return consoleInterface
}

func (ci *ConsoleInterface) ReadLine() (string, error) {
	text, err := ci.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}

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

func (ci *ConsoleInterface) Start(initialPath string) error {
	if initialPath != "" {
		return ci.handleInitialPath(initialPath)
	}

	ci.showMainMenu()
	return nil
}

func (ci *ConsoleInterface) handleInitialPath(pathInput string) error {
	fileInfo, err := os.Stat(pathInput)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("\nPath not found: %s\n", pathInput)
			ci.showMainMenu()
			return nil
		}
		return fmt.Errorf("error checking path: %w", err)
	}

	if fileInfo.IsDir() {
		return ci.handleDirectory(pathInput)
	}

	return ci.handleSingleFile(pathInput)
}

func (ci *ConsoleInterface) showMainMenu() {
	fmt.Print("\033[H\033[2J")
	fmt.Println("====================================")
	fmt.Println("=== DOCUMENT CLASSIFIER ===")
	fmt.Println("====================================")
	fmt.Printf("\nClassifier: %s\n", ci.processingService.GetClassifierName())
	fmt.Printf("Supported formats: %s\n", strings.Join(ci.processingService.GetSupportedFormats(), ", "))

	analyzeService := ci.processingService.GetAnalyzeService()
	fmt.Printf("Rules file: %s\n", analyzeService.GetRulesFilePath())

	fmt.Println("\nSelect an option:")
	fmt.Println("1. Process a single file")
	fmt.Println("2. Process all files in a folder")
	fmt.Println("3. Show current classification rules")
	fmt.Println("4. Reload classification rules")
	fmt.Println("5. Select rules file")
	fmt.Println("6. Exit")
	fmt.Println()

	fmt.Print("Enter your choice (1-6): ")
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
		fmt.Println("Exiting program...")
		os.Exit(0)
	default:
		fmt.Println("\nInvalid option! Press Enter to continue...")
		ci.ReadLine()
		ci.showMainMenu()
	}
}

func (ci *ConsoleInterface) showClassificationRules() {
	fmt.Print("\033[H\033[2J")
	fmt.Println("=== Current Classification Rules ===")

	analyzeService := ci.processingService.GetAnalyzeService()
	rules := analyzeService.GetRules()

	fmt.Printf("Rules file: %s\n\n", analyzeService.GetRulesFilePath())

	for i, rule := range rules {
		fmt.Printf("%d. Type: %s\n", i+1, rule.Type)
		fmt.Printf("   Keywords: %s\n\n", strings.Join(rule.Keywords, ", "))
	}

	fmt.Print("\nPress Enter to return to main menu...")
	ci.ReadLine()
	ci.showMainMenu()
}

func (ci *ConsoleInterface) reloadClassificationRules() {
	fmt.Print("\033[H\033[2J")
	fmt.Println("=== Reloading Classification Rules ===")

	analyzeService := ci.processingService.GetAnalyzeService()
	err := analyzeService.ReloadRules()

	if err != nil {
		fmt.Printf("\nError reloading rules: %v\n", err)
	} else {
		fmt.Printf("\nRules successfully reloaded from file:\n%s\n", analyzeService.GetRulesFilePath())
	}

	fmt.Print("\nPress Enter to return to main menu...")
	ci.ReadLine()
	ci.showMainMenu()
}

func (ci *ConsoleInterface) selectRulesFile() {
	startDir, err := os.Getwd()
	if err != nil {
		startDir = "/"
	}

	fmt.Print("\033[H\033[2J")
	fmt.Println("=== Select JSON Rules File ===")
	fmt.Println("Navigate to the JSON file containing your classification rules.")
	fmt.Print("\nPress Enter to continue...")
	ci.ReadLine()

	selectedPath, err := ci.fileBrowser.BrowseFilesWithFilter(startDir, true, []string{".json"})
	if err != nil {
		fmt.Printf("\nError browsing files: %v\n", err)
		fmt.Print("\nPress Enter to return to main menu...")
		ci.ReadLine()
		ci.showMainMenu()
		return
	}

	if selectedPath == "" {
		ci.showMainMenu()
		return
	}

	analyzeService := ci.processingService.GetAnalyzeService()

	newRules, err := models.LoadRulesFromJSON(selectedPath)
	if err != nil {
		fmt.Printf("\nError loading rules from file: %v\n", err)
		fmt.Print("\nPress Enter to return to main menu...")
		ci.ReadLine()
		ci.showMainMenu()
		return
	}

	err = analyzeService.SetRulesFile(selectedPath)
	if err != nil {
		fmt.Printf("\nError setting new rules file: %v\n", err)
	} else {
		fmt.Printf("\nRules file updated to: %s\n", selectedPath)
		fmt.Printf("Loaded %d classification rules\n", len(newRules))
	}

	fmt.Print("\nPress Enter to return to main menu...")
	ci.ReadLine()
	ci.showMainMenu()
}

func (ci *ConsoleInterface) selectFile() {
	startDir, err := os.Getwd()
	if err != nil {
		startDir = "/"
	}

	selectedPath, err := ci.fileBrowser.BrowseFiles(startDir, true)
	if err != nil {
		fmt.Printf("\nError browsing files: %v\n", err)
		fmt.Print("\nPress Enter to return to main menu...")
		ci.ReadLine()
		ci.showMainMenu()
		return
	}

	if selectedPath == "" {
		ci.showMainMenu()
		return
	}

	if err := ci.handleSingleFile(selectedPath); err != nil {
		fmt.Printf("\nError: %v\n", err)
	}

	fmt.Print("\nPress Enter to return to main menu...")
	ci.ReadLine()
	ci.showMainMenu()
}

func (ci *ConsoleInterface) selectDirectory() {
	startDir, err := os.Getwd()
	if err != nil {
		startDir = "/"
	}

	selectedPath, err := ci.fileBrowser.BrowseFiles(startDir, false)
	if err != nil {
		fmt.Printf("\nError browsing directories: %v\n", err)
		fmt.Print("\nPress Enter to return to main menu...")
		ci.ReadLine()
		ci.showMainMenu()
		return
	}

	if selectedPath == "" {
		ci.showMainMenu()
		return
	}

	if err := ci.handleDirectory(selectedPath); err != nil {
		fmt.Printf("\nError: %v\n", err)
	}

	fmt.Print("\nPress Enter to return to main menu...")
	ci.ReadLine()
	ci.showMainMenu()
}

func (ci *ConsoleInterface) handleSingleFile(filePath string) error {
	fmt.Printf("\nProcessing file: %s\n", filePath)

	document, destinationPath, err := ci.processingService.ProcessSingleFile(filePath)
	if err != nil {
		return err
	}

	fmt.Println("\n===== Classification Result =====")
	fmt.Printf("File: %s\n", document.Filename)
	if document.Classification != nil {
		fmt.Printf("Document type: %s\n", document.Classification.DocumentType)
		fmt.Printf("Keywords: %s\n", strings.Join(document.Classification.Keywords, ", "))
	} else {
		fmt.Println("Could not classify document")
	}
	fmt.Printf("\nFile organized at: %s\n", destinationPath)

	return nil
}

func (ci *ConsoleInterface) handleDirectory(dirPath string) error {
	fmt.Printf("\nProcessing directory: %s\n", dirPath)

	result, err := ci.processingService.ProcessDirectory(dirPath)
	if err != nil {
		return err
	}

	fmt.Println("\n===== Processing Result =====")
	fmt.Printf("Total files processed: %d\n", result.ProcessedCount)
	fmt.Printf("Total failures: %d\n", result.FailedCount)

	for _, fileResult := range result.Results {
		if fileResult.Success {
			fmt.Printf("\n✓ %s → %s\n",
				fileResult.Filename,
				fileResult.DocumentType)
		} else {
			fmt.Printf("\n✗ %s → FAILED: %s\n",
				fileResult.Filename,
				fileResult.Error)
		}
	}

	fmt.Printf("\nDocuments organized at: %s\n", ci.processingService.GetOutputDirectory())

	return nil
}
