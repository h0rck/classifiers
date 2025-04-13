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

// ImageExtractor é um extrator que processa imagens usando OCR
type ImageExtractor struct{}

// ExtractText extrai texto de uma imagem usando OCR
func (e *ImageExtractor) ExtractText(filePath string) (models.DocumentMetadata, error) {
	// Determinar o sistema operacional
	osType := runtime.GOOS

	// Tentar encontrar o Tesseract com base no sistema operacional
	tesseractPath, tesseractErr := e.findTesseract(osType)
	if tesseractErr != nil {
		// Instruções específicas para cada sistema
		installInstructions := e.getInstallInstructions(osType)

		fmt.Println("AVISO: Tesseract OCR não encontrado. " + installInstructions)

		return models.DocumentMetadata{
			Filename: filepath.Base(filePath),
			Text:     fmt.Sprintf("Imagem: %s (OCR não disponível)\n\n%s", filepath.Base(filePath), installInstructions),
		}, nil
	}

	// Criar um arquivo temporário para a saída do tesseract
	tempDir, err := os.MkdirTemp("", "tesseract-output")
	if err != nil {
		return models.DocumentMetadata{}, fmt.Errorf("falha ao criar diretório temporário: %w", err)
	}
	defer os.RemoveAll(tempDir)

	outputPrefix := filepath.Join(tempDir, "output")

	// Executar o tesseract (com suporte para português)
	cmd := exec.Command(tesseractPath, filePath, outputPrefix, "-l", "por")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Verificar se o erro é relacionado ao idioma português
		if strings.Contains(string(output), "Error opening data file") ||
			strings.Contains(string(output), "Failed loading language 'por'") {
			// Tentar com o idioma padrão (geralmente inglês)
			fmt.Println("AVISO: Pacote de idioma português não encontrado. Usando idioma padrão Inglês.")
			cmd = exec.Command(tesseractPath, filePath, outputPrefix)
			output, err = cmd.CombinedOutput()
			if err != nil {
				return models.DocumentMetadata{}, fmt.Errorf("falha ao executar OCR: %w\nSaída: %s", err, string(output))
			}
		} else {
			return models.DocumentMetadata{}, fmt.Errorf("falha ao executar OCR: %w\nSaída: %s", err, string(output))
		}
	}

	// Ler o arquivo de saída
	textBytes, err := os.ReadFile(outputPrefix + ".txt")
	if err != nil {
		return models.DocumentMetadata{}, fmt.Errorf("falha ao ler saída do OCR: %w", err)
	}

	text := e.postprocessText(string(textBytes))

	return models.DocumentMetadata{
		Filename: filepath.Base(filePath),
		Text:     text,
	}, nil
}

// findTesseract tenta encontrar o executável do Tesseract com base no sistema operacional
func (e *ImageExtractor) findTesseract(osType string) (string, error) {
	// Primeiro, tentar encontrar no PATH (funciona em qualquer sistema operacional)
	tesseractPath, err := exec.LookPath("tesseract")
	if err == nil {
		return tesseractPath, nil
	}

	// Verificar caminhos comuns com base no sistema operacional
	var possiblePaths []string

	switch osType {
	case "windows":
		possiblePaths = []string{
			"C:\\Program Files\\Tesseract-OCR\\tesseract.exe",
			"C:\\Program Files (x86)\\Tesseract-OCR\\tesseract.exe",
			"C:\\Tesseract-OCR\\tesseract.exe",
		}
	case "darwin": // macOS
		possiblePaths = []string{
			"/usr/local/bin/tesseract",
			"/opt/homebrew/bin/tesseract",
			"/opt/local/bin/tesseract", // MacPorts
		}
	case "linux":
		possiblePaths = []string{
			"/usr/bin/tesseract",
			"/usr/local/bin/tesseract",
			"/snap/bin/tesseract",
		}
	}

	// Verificar cada caminho possível
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("tesseract não encontrado")
}

// getInstallInstructions retorna instruções específicas para o sistema operacional
func (e *ImageExtractor) getInstallInstructions(osType string) string {
	switch osType {
	case "windows":
		return "Para instalar o Tesseract OCR no Windows:\n" +
			"1. Acesse https://github.com/UB-Mannheim/tesseract/wiki\n" +
			"2. Baixe o instalador mais recente (tesseract-ocr-w64-setup-v5.x.x.exe)\n" +
			"3. Durante a instalação, selecione o idioma português\n" +
			"4. Reinicie o aplicativo após a instalação"

	case "darwin": // macOS
		return "Para instalar o Tesseract OCR no macOS:\n" +
			"1. Instale o Homebrew se ainda não tiver (https://brew.sh/)\n" +
			"2. Execute no terminal: brew install tesseract tesseract-lang\n" +
			"3. Reinicie o aplicativo após a instalação"

	case "linux":
		return "Para instalar o Tesseract OCR no Ubuntu/Debian:\n" +
			"sudo apt-get update && sudo apt-get install -y tesseract-ocr tesseract-ocr-por\n\n" +
			"Para Fedora/RHEL:\n" +
			"sudo dnf install -y tesseract tesseract-langpack-por"

	default:
		return "Para instalar o Tesseract OCR, consulte as instruções em: https://github.com/tesseract-ocr/tesseract"
	}
}

// postprocessText limpa e formata o texto extraído pela OCR
func (e *ImageExtractor) postprocessText(text string) string {
	// Se o texto estiver vazio, retornar mensagem
	if strings.TrimSpace(text) == "" {
		return "Não foi possível extrair texto desta imagem."
	}

	// Remover espaços extras
	re := regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, " ")

	// Substituir quebras de linha por espaços
	text = strings.ReplaceAll(text, "\r\n", " ")
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\r", " ")

	return strings.TrimSpace(text)
}

// IsSupportedFormat verifica se o formato do arquivo é suportado
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

// GetSupportedFormats retorna os formatos suportados por este extrator
func (e *ImageExtractor) GetSupportedFormats() []string {
	return []string{".png", ".jpg", ".jpeg", ".bmp", ".tiff", ".tif", ".gif"}
}
