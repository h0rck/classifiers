package extractors

import (
	"fmt"
	"path/filepath"
	"relatorios/interfaces"
)

type DocumentExtractorFactory struct {
	extractors []interfaces.TextExtractor
}

func NewDocumentExtractorFactory() *DocumentExtractorFactory {
	return &DocumentExtractorFactory{
		extractors: []interfaces.TextExtractor{
			&PdfExtractor{},
			&WordExtractor{},
			&ExcelExtractor{},
			&TextFileExtractor{},
			&ImageExtractor{},
		},
	}
}

func (f *DocumentExtractorFactory) GetExtractorForFile(filePath string) (interfaces.TextExtractor, error) {
	for _, extractor := range f.extractors {
		if extractor.IsSupportedFormat(filePath) {
			return extractor, nil
		}
	}

	return nil, fmt.Errorf("unsupported file format: %s", filepath.Ext(filePath))
}

func (f *DocumentExtractorFactory) IsFormatSupported(filePath string) bool {
	for _, extractor := range f.extractors {
		if extractor.IsSupportedFormat(filePath) {
			return true
		}
	}
	return false
}

func (f *DocumentExtractorFactory) GetSupportedFormats() []string {
	var formats []string
	for _, extractor := range f.extractors {
		formats = append(formats, extractor.GetSupportedFormats()...)
	}
	return formats
}
