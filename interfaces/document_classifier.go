package interfaces

import "relatorios/models"

// DocumentClassifier define a interface para classificadores de documentos
type DocumentClassifier interface {
	// Classify analisa e classifica um documento com base em seu conteúdo
	Classify(document models.DocumentMetadata) (models.DocumentMetadata, error)

	// GetClassifierName retorna o nome do classificador
	GetClassifierName() string
}
