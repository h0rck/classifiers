package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type FileBrowser struct {
	reader        InputReader
	formatChecker FormatChecker
}

type InputReader interface {
	ReadLine() (string, error)
}

type FormatChecker interface {
	IsFormatSupported(filename string) bool
}

func NewFileBrowser(reader InputReader, formatChecker FormatChecker) *FileBrowser {
	return &FileBrowser{
		reader:        reader,
		formatChecker: formatChecker,
	}
}

func (fb *FileBrowser) BrowseFiles(currentDir string, selectFile bool) (string, error) {
	return fb.BrowseFilesWithFilter(currentDir, selectFile, nil)
}

func (fb *FileBrowser) BrowseFilesWithFilter(currentDir string, selectFile bool, allowedExtensions []string) (string, error) {
	for {
		fmt.Print("\033[H\033[2J")

		fmt.Printf("Current directory: %s\n\n", currentDir)

		if selectFile {
			if len(allowedExtensions) > 0 {
				fmt.Printf("=== Select a %s file or navigate through directories ===\n", strings.Join(allowedExtensions, "/"))
			} else {
				fmt.Println("=== Select a FILE or navigate through directories ===")
			}
		} else {
			fmt.Println("=== Select a DIRECTORY ===")
		}

		fmt.Println("[0] .. (Go back)")
		fmt.Println("[C] Cancel")
		fmt.Println("[M] Enter path manually")
		if !selectFile {
			fmt.Println("[S] Select current directory")
		}
		fmt.Println("-------------------------------------")

		entries, err := os.ReadDir(currentDir)
		if err != nil {
			return "", err
		}

		var dirs []os.DirEntry
		var files []os.DirEntry

		for _, entry := range entries {
			if entry.IsDir() {
				dirs = append(dirs, entry)
			} else if selectFile {
				if len(allowedExtensions) > 0 {
					ext := strings.ToLower(filepath.Ext(entry.Name()))
					for _, allowedExt := range allowedExtensions {
						if ext == allowedExt {
							files = append(files, entry)
							break
						}
					}
				} else if fb.formatChecker.IsFormatSupported(entry.Name()) {
					files = append(files, entry)
				}
			}
		}

		sort.Slice(dirs, func(i, j int) bool {
			return dirs[i].Name() < dirs[j].Name()
		})
		sort.Slice(files, func(i, j int) bool {
			return files[i].Name() < files[j].Name()
		})

		for i, entry := range dirs {
			fmt.Printf("[%d] 📁 %s\n", i+1, entry.Name())
		}

		if selectFile {
			for i, entry := range files {
				fmt.Printf("[%d] 📄 %s\n", i+len(dirs)+1, entry.Name())
			}
		}

		fmt.Print("\nSelect an option (number, C, M" + func() string {
			if !selectFile {
				return ", S"
			}
			return ""
		}() + "): ")

		input, _ := fb.reader.ReadLine()
		input = strings.ToUpper(strings.TrimSpace(input))

		switch input {
		case "0":
			parent := filepath.Dir(currentDir)
			if parent == currentDir {
				continue
			}
			currentDir = parent
			continue
		case "C":
			return "", nil
		case "M":
			fmt.Print("\nEnter the complete path: ")
			path, _ := fb.reader.ReadLine()
			path = strings.TrimSpace(path)

			if path == "" {
				continue
			}

			info, err := os.Stat(path)
			if err != nil {
				fmt.Printf("\nError: Invalid or inaccessible path.\n")
				fmt.Print("Press Enter to continue...")
				fb.reader.ReadLine()
				continue
			}

			if selectFile && info.IsDir() {
				currentDir = path
				continue
			}

			if !selectFile && !info.IsDir() {
				fmt.Printf("\nError: The selected path is not a directory.\n")
				fmt.Print("Press Enter to continue...")
				fb.reader.ReadLine()
				continue
			}

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
					fmt.Printf("\nError: The file must have one of the following extensions: %s\n",
						strings.Join(allowedExtensions, ", "))
					fmt.Print("Press Enter to continue...")
					fb.reader.ReadLine()
					continue
				}
			}

			return path, nil
		case "S":
			if !selectFile {
				return currentDir, nil
			}
		}

		var index int
		if _, err := fmt.Sscanf(input, "%d", &index); err != nil {
			continue
		}

		if index >= 1 && index <= len(dirs) {
			currentDir = filepath.Join(currentDir, dirs[index-1].Name())
			continue
		}

		if selectFile && index > len(dirs) && index <= len(dirs)+len(files) {
			fileIndex := index - len(dirs) - 1
			selectedFile := filepath.Join(currentDir, files[fileIndex].Name())
			return selectedFile, nil
		}
	}
}
