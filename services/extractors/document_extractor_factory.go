package extractors

import (
	"fmt"
	"path/filepath"
	"relatorios/interfaces"
)

// DocumentExtractorFactory é responsável por criar os extratores apropriados
type DocumentExtractorFactory struct {
	extractors []interfaces.TextExtractor
}

// NewDocumentExtractorFactory cria uma nova instância da fábrica de extratores
func NewDocumentExtractorFactory() *DocumentExtractorFactory {
	return &DocumentExtractorFactory{
		extractors: []interfaces.TextExtractor{
			&PdfExtractor{},
			&WordExtractor{},
			&ExcelExtractor{},
			&TextFileExtractor{},
		},
	}
}

// GetExtractorForFile retorna o extrator apropriado para o tipo de arquivo
func (f *DocumentExtractorFactory) GetExtractorForFile(filePath string) (interfaces.TextExtractor, error) {
	for _, extractor := range f.extractors {
		if extractor.IsSupportedFormat(filePath) {
			return extractor, nil
		}
	}

	return nil, fmt.Errorf("formato de arquivo não suportado: %s", filepath.Ext(filePath))
}

// IsFormatSupported verifica se existe um extrator para o formato do arquivo
func (f *DocumentExtractorFactory) IsFormatSupported(filePath string) bool {
	for _, extractor := range f.extractors {
		if extractor.IsSupportedFormat(filePath) {
			return true
		}
	}
	return false
}

// GetSupportedFormats retorna todos os formatos suportados por todos os extratores
func (f *DocumentExtractorFactory) GetSupportedFormats() []string {
	var formats []string
	for _, extractor := range f.extractors {
		formats = append(formats, extractor.GetSupportedFormats()...)
	}
	return formats
}
