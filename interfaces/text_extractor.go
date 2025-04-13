package interfaces

import "relatorios/models"

// TextExtractor define a interface para extrair texto de documentos
type TextExtractor interface {
	ExtractText(filePath string) (models.DocumentMetadata, error)
	IsSupportedFormat(filePath string) bool
	GetSupportedFormats() []string
}
