package extractors

import (
	"os"
	"path/filepath"
	"relatorios/models"
)

// TextFileExtractor é um extrator que lê arquivos de texto simples
type TextFileExtractor struct{}

// ExtractText extrai o texto de um arquivo .txt
func (e *TextFileExtractor) ExtractText(filePath string) (models.DocumentMetadata, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return models.DocumentMetadata{}, err
	}

	return models.DocumentMetadata{
		Filename: filepath.Base(filePath),
		Text:     string(data),
	}, nil
}

// IsSupportedFormat verifica se o formato do arquivo é suportado
func (e *TextFileExtractor) IsSupportedFormat(filePath string) bool {
	ext := filepath.Ext(filePath)
	return ext == ".txt"
}

// GetSupportedFormats retorna os formatos suportados por este extrator
func (e *TextFileExtractor) GetSupportedFormats() []string {
	return []string{".txt"}
}
