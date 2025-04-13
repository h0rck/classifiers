package extractors

import (
	"fmt"
	"path/filepath"
	"relatorios/models"

	"github.com/ledongthuc/pdf"
)

type PdfExtractor struct{}

func (e *PdfExtractor) ExtractText(filePath string) (models.DocumentMetadata, error) {
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return models.DocumentMetadata{}, fmt.Errorf("failed to open PDF: %w", err)
	}
	defer f.Close()

	var text string
	totalPage := r.NumPage()

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		pageText, err := p.GetPlainText(nil)
		if err != nil {
			return models.DocumentMetadata{}, fmt.Errorf("failed to extract text from page %d: %w", pageIndex, err)
		}
		text += pageText
	}

	return models.DocumentMetadata{
		Filename: filepath.Base(filePath),
		Text:     text,
	}, nil
}

func (e *PdfExtractor) IsSupportedFormat(filePath string) bool {
	ext := filepath.Ext(filePath)
	return ext == ".pdf"
}

func (e *PdfExtractor) GetSupportedFormats() []string {
	return []string{".pdf"}
}
