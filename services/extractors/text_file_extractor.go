package extractors

import (
	"os"
	"path/filepath"
	"relatorios/models"
)

type TextFileExtractor struct{}

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

func (e *TextFileExtractor) IsSupportedFormat(filePath string) bool {
	ext := filepath.Ext(filePath)
	return ext == ".txt"
}

func (e *TextFileExtractor) GetSupportedFormats() []string {
	return []string{".txt"}
}
