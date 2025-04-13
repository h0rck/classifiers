package classifiers

import (
	"relatorios/interfaces"
	"relatorios/models"
)

type DocumentClassifier struct {
	analyzeService interfaces.AnalyzeService
}

func NewDocumentClassifier(analyzeService interfaces.AnalyzeService) *DocumentClassifier {
	return &DocumentClassifier{
		analyzeService: analyzeService,
	}
}

func (c *DocumentClassifier) Classify(document models.DocumentMetadata) (models.DocumentMetadata, error) {
	result := c.analyzeService.Execute(document.Text)
	document.Classification = &result.Classification
	return document, nil
}

func (c *DocumentClassifier) GetClassifierName() string {
	return "Keyword Classifier"
}

func (c *DocumentClassifier) GetAnalyzeService() interfaces.AnalyzeService {
	return c.analyzeService
}
