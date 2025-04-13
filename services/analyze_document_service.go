package services

import (
	"regexp"
	"relatorios/models"
	"strings"
)

// AnalyzeDocumentService é responsável por analisar e classificar documentos
type AnalyzeDocumentService struct {
	rules []models.DocumentRule
}

// NewAnalyzeDocumentService cria uma nova instância do serviço
func NewAnalyzeDocumentService() *AnalyzeDocumentService {
	return &AnalyzeDocumentService{
		rules: models.GetDefaultRules(),
	}
}

// SetRules permite alterar as regras de classificação
func (s *AnalyzeDocumentService) SetRules(rules []models.DocumentRule) {
	s.rules = rules
}

// Execute analisa o texto do documento e gera uma classificação
func (s *AnalyzeDocumentService) Execute(text string) *models.ClassificationResult {
	// Se não tiver texto, retorna como documento vazio
	if text == "" {
		return s.createResult("Documento Vazio", []string{"vazio"})
	}

	// Normalizar o texto para comparação
	normalizedText := s.normalizeText(text)

	// Verificar cada regra e encontrar a que tem mais palavras-chave
	bestMatchCount := 0
	bestType := ""
	bestKeywords := []string{}

	for _, rule := range s.rules {
		matchedKeywords := []string{}

		// Verificar cada palavra-chave
		for _, keyword := range rule.Keywords {
			normalizedKeyword := strings.ToLower(keyword)
			if strings.Contains(normalizedText, normalizedKeyword) {
				matchedKeywords = append(matchedKeywords, keyword)
			}
		}

		// Atualizar a melhor correspondência
		if len(matchedKeywords) > bestMatchCount {
			bestMatchCount = len(matchedKeywords)
			bestType = rule.Type
			bestKeywords = matchedKeywords
		}
	}

	// Se encontrou uma classificação
	if bestType != "" {
		// Limitar a 5 palavras-chave no máximo
		if len(bestKeywords) > 5 {
			bestKeywords = bestKeywords[:5]
		}

		return s.createResult(bestType, bestKeywords)
	}

	// Se não encontramos nenhuma correspondência, retornamos como "Outro"
	return s.createResult("Outro", []string{"documento", "texto"})
}

// createResult cria um objeto de resultado de classificação
func (s *AnalyzeDocumentService) createResult(documentType string, keywords []string) *models.ClassificationResult {
	return &models.ClassificationResult{
		Classification: models.DocumentClassification{
			DocumentType: documentType,
			Keywords:     keywords,
		},
	}
}

// normalizeText normaliza o texto para facilitar a comparação
func (s *AnalyzeDocumentService) normalizeText(text string) string {
	// Converter para minúsculas
	normalized := strings.ToLower(strings.TrimSpace(text))

	// Substituir múltiplos espaços em branco por um único espaço
	re := regexp.MustCompile(`\s+`)
	normalized = re.ReplaceAllString(normalized, " ")

	return normalized
}
