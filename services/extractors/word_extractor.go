package extractors

import (
	"fmt"
	"path/filepath"
	"relatorios/models"

	"github.com/unidoc/unioffice/document"
)

// WordExtractor é um extrator que lê arquivos DOCX
type WordExtractor struct{}

// ExtractText extrai o texto de um arquivo Word (DOCX)
func (e *WordExtractor) ExtractText(filePath string) (models.DocumentMetadata, error) {
	doc, err := document.Open(filePath)
	if err != nil {
		return models.DocumentMetadata{}, fmt.Errorf("falha ao abrir documento Word: %w", err)
	}

	var text string
	for _, para := range doc.Paragraphs() {
		for _, run := range para.Runs() {
			text += run.Text()
		}
		text += "\n"
	}

	return models.DocumentMetadata{
		Filename: filepath.Base(filePath),
		Text:     text,
	}, nil
}

// IsSupportedFormat verifica se o formato do arquivo é suportado
func (e *WordExtractor) IsSupportedFormat(filePath string) bool {
	ext := filepath.Ext(filePath)
	return ext == ".docx" || ext == ".doc"
}

// GetSupportedFormats retorna os formatos suportados por este extrator
func (e *WordExtractor) GetSupportedFormats() []string {
	return []string{".docx", ".doc"}
}
