package extractors

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"relatorios/models"
)

type ImageExtractor struct{}

func (e *ImageExtractor) ExtractText(filePath string) (models.DocumentMetadata, error) {
	osType := runtime.GOOS

	tesseractPath, tesseractErr := e.findTesseract(osType)
	if tesseractErr != nil {
		installInstructions := e.getInstallInstructions(osType)

		fmt.Println("WARNING: Tesseract OCR not found. " + installInstructions)

		return models.DocumentMetadata{
			Filename: filepath.Base(filePath),
			Text:     fmt.Sprintf("Image: %s (OCR not available)\n\n%s", filepath.Base(filePath), installInstructions),
		}, nil
	}

	tempDir, err := os.MkdirTemp("", "tesseract-output")
	if err != nil {
		return models.DocumentMetadata{}, fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	outputPrefix := filepath.Join(tempDir, "output")

	cmd := exec.Command(tesseractPath, filePath, outputPrefix, "-l", "por")
	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(output), "Error opening data file") ||
			strings.Contains(string(output), "Failed loading language 'por'") {
			fmt.Println("WARNING: Portuguese language package not found. Using default English.")
			cmd = exec.Command(tesseractPath, filePath, outputPrefix)
			output, err = cmd.CombinedOutput()
			if err != nil {
				return models.DocumentMetadata{}, fmt.Errorf("failed to execute OCR: %w\nOutput: %s", err, string(output))
			}
		} else {
			return models.DocumentMetadata{}, fmt.Errorf("failed to execute OCR: %w\nOutput: %s", err, string(output))
		}
	}

	textBytes, err := os.ReadFile(outputPrefix + ".txt")
	if err != nil {
		return models.DocumentMetadata{}, fmt.Errorf("failed to read OCR output: %w", err)
	}

	text := e.postprocessText(string(textBytes))

	return models.DocumentMetadata{
		Filename: filepath.Base(filePath),
		Text:     text,
	}, nil
}

func (e *ImageExtractor) findTesseract(osType string) (string, error) {
	tesseractPath, err := exec.LookPath("tesseract")
	if err == nil {
		return tesseractPath, nil
	}

	var possiblePaths []string

	switch osType {
	case "windows":
		possiblePaths = []string{
			"C:\\Program Files\\Tesseract-OCR\\tesseract.exe",
			"C:\\Program Files (x86)\\Tesseract-OCR\\tesseract.exe",
			"C:\\Tesseract-OCR\\tesseract.exe",
		}
	case "darwin":
		possiblePaths = []string{
			"/usr/local/bin/tesseract",
			"/opt/homebrew/bin/tesseract",
			"/opt/local/bin/tesseract",
		}
	case "linux":
		possiblePaths = []string{
			"/usr/bin/tesseract",
			"/usr/local/bin/tesseract",
			"/snap/bin/tesseract",
		}
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("tesseract not found")
}

func (e *ImageExtractor) getInstallInstructions(osType string) string {
	switch osType {
	case "windows":
		return "To install Tesseract OCR on Windows:\n" +
			"1. Visit https://github.com/UB-Mannheim/tesseract/wiki\n" +
			"2. Download the latest installer (tesseract-ocr-w64-setup-v5.x.x.exe)\n" +
			"3. During installation, select desired languages\n" +
			"4. Restart the application after installation"

	case "darwin":
		return "To install Tesseract OCR on macOS:\n" +
			"1. Install Homebrew if you don't have it (https://brew.sh/)\n" +
			"2. Run in terminal: brew install tesseract tesseract-lang\n" +
			"3. Restart the application after installation"

	case "linux":
		return "To install Tesseract OCR on Ubuntu/Debian:\n" +
			"sudo apt-get update && sudo apt-get install -y tesseract-ocr tesseract-ocr-eng\n\n" +
			"For Fedora/RHEL:\n" +
			"sudo dnf install -y tesseract tesseract-langpack-eng"

	default:
		return "To install Tesseract OCR, see instructions at: https://github.com/tesseract-ocr/tesseract"
	}
}

func (e *ImageExtractor) postprocessText(text string) string {
	if strings.TrimSpace(text) == "" {
		return "Could not extract text from this image."
	}

	re := regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, " ")

	text = strings.ReplaceAll(text, "\r\n", " ")
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\r", " ")

	return strings.TrimSpace(text)
}

func (e *ImageExtractor) IsSupportedFormat(filePath string) bool {
	supportedFormats := []string{".png", ".jpg", ".jpeg", ".bmp", ".tiff", ".tif", ".gif"}
	ext := strings.ToLower(filepath.Ext(filePath))

	for _, format := range supportedFormats {
		if ext == format {
			return true
		}
	}
	return false
}

func (e *ImageExtractor) GetSupportedFormats() []string {
	return []string{".png", ".jpg", ".jpeg", ".bmp", ".tiff", ".tif", ".gif"}
}
