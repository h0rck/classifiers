package interfaces

import "relatorios/models"

// DocumentClassifier define a interface para classificadores de documentos
type DocumentClassifier interface {
	Classify(document models.DocumentMetadata) (models.DocumentMetadata, error)
	GetClassifierName() string
}
