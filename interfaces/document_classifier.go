package interfaces

import "relatorios/models"

type DocumentClassifier interface {
	Classify(document models.DocumentMetadata) (models.DocumentMetadata, error)
	GetClassifierName() string
}
