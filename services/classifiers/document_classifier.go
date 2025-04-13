package classifiers

import (
	"relatorios/models"
	"relatorios/services"
)

// DocumentClassifier implementa a classificação de documentos
type DocumentClassifier struct {
	analyzeDocument *services.AnalyzeDocumentService
}

// NewDocumentClassifier cria uma nova instância do classificador
func NewDocumentClassifier(analyzeDocument *services.AnalyzeDocumentService) *DocumentClassifier {
	return &DocumentClassifier{
		analyzeDocument: analyzeDocument,
	}
}

// GetClassifierName retorna o nome do classificador
func (c *DocumentClassifier) GetClassifierName() string {
	return "Classificador por script"
}

// Classify classifica um documento com base em seu conteúdo
func (c *DocumentClassifier) Classify(document models.DocumentMetadata) (models.DocumentMetadata, error) {
	result := c.analyzeDocument.Execute(document.Text)
	document.Classification = &result.Classification

	return document, nil
}
