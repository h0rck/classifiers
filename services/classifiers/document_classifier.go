package classifiers

import (
	"relatorios/interfaces"
	"relatorios/models"
)

// DocumentClassifier implementa a interface DocumentClassifier
type DocumentClassifier struct {
	analyzeService interfaces.AnalyzeService
}

// NewDocumentClassifier cria uma nova instância do classificador
func NewDocumentClassifier(analyzeService interfaces.AnalyzeService) *DocumentClassifier {
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

// GetAnalyzeService retorna o serviço de análise usado pelo classificador
func (c *DocumentClassifier) GetAnalyzeService() interfaces.AnalyzeService {
	return c.analyzeService
}
