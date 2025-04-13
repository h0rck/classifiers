package extractors

import (
	"fmt"
	"path/filepath"
	"relatorios/models"

	"github.com/unidoc/unioffice/document"
)

type WordExtractor struct{}

func (e *WordExtractor) ExtractText(filePath string) (models.DocumentMetadata, error) {
	doc, err := document.Open(filePath)
	if err != nil {
		return models.DocumentMetadata{}, fmt.Errorf("failed to open Word document: %w", err)
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

func (e *WordExtractor) IsSupportedFormat(filePath string) bool {
	ext := filepath.Ext(filePath)
	return ext == ".docx" || ext == ".doc"
}

func (e *WordExtractor) GetSupportedFormats() []string {
	return []string{".docx", ".doc"}
}
