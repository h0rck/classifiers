package extractors

import (
	"fmt"
	"path/filepath"
	"relatorios/models"

	"github.com/ledongthuc/pdf"
)

// PdfExtractor é um extrator que lê arquivos PDF
type PdfExtractor struct{}

// ExtractText extrai o texto de um arquivo PDF
func (e *PdfExtractor) ExtractText(filePath string) (models.DocumentMetadata, error) {
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return models.DocumentMetadata{}, fmt.Errorf("falha ao abrir PDF: %w", err)
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
			return models.DocumentMetadata{}, fmt.Errorf("falha ao extrair texto da página %d: %w", pageIndex, err)
		}
		text += pageText
	}

	return models.DocumentMetadata{
		Filename: filepath.Base(filePath),
		Text:     text,
	}, nil
}

// IsSupportedFormat verifica se o formato do arquivo é suportado
func (e *PdfExtractor) IsSupportedFormat(filePath string) bool {
	ext := filepath.Ext(filePath)
	return ext == ".pdf"
}

// GetSupportedFormats retorna os formatos suportados por este extrator
func (e *PdfExtractor) GetSupportedFormats() []string {
	return []string{".pdf"}
}
