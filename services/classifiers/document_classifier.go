package classifiers

import (
	"relatorios/models"
	"relatorios/services"
)

// DocumentClassifier implementa a interface DocumentClassifier
type DocumentClassifier struct {
	analyzeService *services.AnalyzeDocumentService
}

// NewDocumentClassifier cria uma nova instância do classificador
func NewDocumentClassifier(analyzeService *services.AnalyzeDocumentService) *DocumentClassifier {
	return &DocumentClassifier{
		analyzeService: analyzeService,
	}
}

// Classify analisa e classifica um documento
func (c *DocumentClassifier) Classify(document models.DocumentMetadata) (models.DocumentMetadata, error) {
	result := c.analyzeService.Execute(document.Text)
	document.Classification = &result.Classification
	return document, nil
}

// GetClassifierName retorna o nome do classificador
func (c *DocumentClassifier) GetClassifierName() string {
	return "Classificador de Palavras-Chave"
}
