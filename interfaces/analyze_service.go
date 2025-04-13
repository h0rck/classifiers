package interfaces

import "relatorios/models"

// AnalyzeService define a interface para serviços de análise de documentos
type AnalyzeService interface {
	// Execute analisa o texto do documento e gera uma classificação
	Execute(text string) *models.ClassificationResult

	// GetRules retorna as regras de classificação atuais
	GetRules() []models.DocumentRule

	// ReloadRules recarrega as regras do arquivo JSON
	ReloadRules() error

	// GetRulesFilePath retorna o caminho do arquivo JSON de regras
	GetRulesFilePath() string

	// SetRules permite alterar as regras de classificação
	SetRules(rules []models.DocumentRule)
}
