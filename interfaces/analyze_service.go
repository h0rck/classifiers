package interfaces

import "relatorios/models"

type AnalyzeService interface {
	Execute(text string) *models.ClassificationResult
	GetRules() []models.DocumentRule
	ReloadRules() error
	GetRulesFilePath() string
	SetRules(rules []models.DocumentRule)
	SetRulesFile(filePath string) error
}
