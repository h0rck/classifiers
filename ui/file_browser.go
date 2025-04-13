package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// FileBrowser provides file navigation functionality
type FileBrowser struct {
	reader        InputReader
	formatChecker FormatChecker
}

// InputReader defines an interface for reading input
type InputReader interface {
	ReadLine() (string, error)
}

// FormatChecker defines an interface for checking file format support
type FormatChecker interface {
	IsFormatSupported(filename string) bool
}

// NewFileBrowser creates a new FileBrowser instance
func NewFileBrowser(reader InputReader, formatChecker FormatChecker) *FileBrowser {
	return &FileBrowser{
		reader:        reader,
		formatChecker: formatChecker,
	}
}

// BrowseFiles allows the user to navigate through files and directories
func (fb *FileBrowser) BrowseFiles(currentDir string, selectFile bool) (string, error) {
	return fb.BrowseFilesWithFilter(currentDir, selectFile, nil)
}

// BrowseFilesWithFilter allows the user to navigate through files and directories with a custom filter
func (fb *FileBrowser) BrowseFilesWithFilter(currentDir string, selectFile bool, allowedExtensions []string) (string, error) {
	for {
		fmt.Print("\033[H\033[2J") // Clear the screen

		// Show current directory
		fmt.Printf("Diretório atual: %s\n\n", currentDir)

		if selectFile {
			if len(allowedExtensions) > 0 {
				fmt.Printf("=== Selecione um arquivo %s ou navegue pelos diretórios ===\n", strings.Join(allowedExtensions, "/"))
			} else {
				fmt.Println("=== Selecione um ARQUIVO ou navegue pelos diretórios ===")
			}
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

		// Read files and directories
		entries, err := os.ReadDir(currentDir)
		if err != nil {
			return "", err
		}

		// Separate directories and files
		var dirs []os.DirEntry
		var files []os.DirEntry

		// Filter compatible files and directories
		for _, entry := range entries {
			if entry.IsDir() {
				dirs = append(dirs, entry)
			} else if selectFile {
				// If we're selecting a file
				if len(allowedExtensions) > 0 {
					// Check allowed extensions
					ext := strings.ToLower(filepath.Ext(entry.Name()))
					for _, allowedExt := range allowedExtensions {
						if ext == allowedExt {
							files = append(files, entry)
							break
						}
					}
				} else if fb.formatChecker.IsFormatSupported(entry.Name()) {
					// No specific extensions, use default supported formats
					files = append(files, entry)
				}
			}
		}

		// Sort directories and files alphabetically
		sort.Slice(dirs, func(i, j int) bool {
			return dirs[i].Name() < dirs[j].Name()
		})
		sort.Slice(files, func(i, j int) bool {
			return files[i].Name() < files[j].Name()
		})

		// List directories first
		for i, entry := range dirs {
			fmt.Printf("[%d] 📁 %s\n", i+1, entry.Name())
		}

		// Then list files if we're selecting a file
		if selectFile {
			for i, entry := range files {
				fmt.Printf("[%d] 📄 %s\n", i+len(dirs)+1, entry.Name())
			}
		}

		// Request choice
		fmt.Print("\nEscolha uma opção (número, C, M" + func() string {
			if !selectFile {
				return ", S"
			}
			return ""
		}() + "): ")

		input, _ := fb.reader.ReadLine()
		input = strings.ToUpper(strings.TrimSpace(input))

		// Check special options
		switch input {
		case "0": // Go back
			parent := filepath.Dir(currentDir)
			if parent == currentDir { // We're already at the root
				continue
			}
			currentDir = parent
			continue
		case "C": // Cancel
			return "", nil
		case "M": // Enter path manually
			fmt.Print("\nDigite o caminho completo: ")
			path, _ := fb.reader.ReadLine()
			path = strings.TrimSpace(path)

			if path == "" {
				continue
			}

			// Check if the path exists
			info, err := os.Stat(path)
			if err != nil {
				fmt.Printf("\nErro: Caminho inválido ou inacessível.\n")
				fmt.Print("Pressione Enter para continuar...")
				fb.reader.ReadLine()
				continue
			}

			// If we want a file but selected a directory, navigate to it
			if selectFile && info.IsDir() {
				currentDir = path
				continue
			}

			// If we want a directory but selected a file, show an error
			if !selectFile && !info.IsDir() {
				fmt.Printf("\nErro: O caminho selecionado não é um diretório.\n")
				fmt.Print("Pressione Enter para continuar...")
				fb.reader.ReadLine()
				continue
			}

			// If we have specific extensions, check if the selected file is valid
			if selectFile && len(allowedExtensions) > 0 && !info.IsDir() {
				ext := strings.ToLower(filepath.Ext(path))
				valid := false
				for _, allowedExt := range allowedExtensions {
					if ext == allowedExt {
						valid = true
						break
					}
				}

				if !valid {
					fmt.Printf("\nErro: O arquivo deve ter uma das seguintes extensões: %s\n",
						strings.Join(allowedExtensions, ", "))
					fmt.Print("Pressione Enter para continuar...")
					fb.reader.ReadLine()
					continue
				}
			}

			return path, nil
		case "S": // Select current directory (only if !selectFile)
			if !selectFile {
				return currentDir, nil
			}
		}

		// Convert choice to number
		var index int
		if _, err := fmt.Sscanf(input, "%d", &index); err != nil {
			continue // If it's not a number, ignore
		}

		// Check if it's a directory
		if index >= 1 && index <= len(dirs) {
			currentDir = filepath.Join(currentDir, dirs[index-1].Name())
			continue
		}

		// Check if it's a file and we're selecting a file
		if selectFile && index > len(dirs) && index <= len(dirs)+len(files) {
			fileIndex := index - len(dirs) - 1
			selectedFile := filepath.Join(currentDir, files[fileIndex].Name())
			return selectedFile, nil
		}
	}
}
