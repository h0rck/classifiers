package interfaces

import "relatorios/models"

type TextExtractor interface {
	ExtractText(filePath string) (models.DocumentMetadata, error)
	IsSupportedFormat(filePath string) bool
	GetSupportedFormats() []string
}
