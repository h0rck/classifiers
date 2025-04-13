package extractors

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"relatorios/models"
)

// ImageExtractor é um extrator que processa imagens usando OCR
type ImageExtractor struct{}

// ExtractText extrai texto de uma imagem usando OCR
func (e *ImageExtractor) ExtractText(filePath string) (models.DocumentMetadata, error) {
	// Procurar o tesseract no PATH do sistema
	tesseractPath, err := exec.LookPath("tesseract")
	if err != nil {
		// Se não encontrar, tentar caminho padrão
		tesseractPath = "/usr/bin/tesseract"
		_, err = os.Stat(tesseractPath)
		if err != nil {
			return models.DocumentMetadata{
				Filename: filepath.Base(filePath),
				Text:     fmt.Sprintf("Imagem: %s (OCR não disponível)", filepath.Base(filePath)),
			}, nil
		}
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
		// Se falhar com português, tente com inglês
		cmd = exec.Command(tesseractPath, filePath, outputPrefix)
		output, err = cmd.CombinedOutput()
		if err != nil {
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

// postprocessText limpa e formata o texto extraído pela OCR
func (e *ImageExtractor) postprocessText(text string) string {
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
