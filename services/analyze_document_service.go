package services

import (
	"fmt"
	"os"
	"regexp"
	"relatorios/models"
	"strings"
)

type AnalyzeDocumentService struct {
	rules     []models.DocumentRule
	rulesFile string
}

func NewAnalyzeDocumentService(rulesFile string) *AnalyzeDocumentService {
	rules, _ := models.LoadRulesFromJSON(rulesFile)

	return &AnalyzeDocumentService{
		rules:     rules,
		rulesFile: rulesFile,
	}
}

func (s *AnalyzeDocumentService) SetRules(rules []models.DocumentRule) {
	s.rules = rules
	_ = models.SaveRulesToJSON(s.rulesFile, rules)
}

func (s *AnalyzeDocumentService) GetRules() []models.DocumentRule {
	return s.rules
}

func (s *AnalyzeDocumentService) ReloadRules() error {
	rules, err := models.LoadRulesFromJSON(s.rulesFile)
	if err != nil {
		return err
	}
	s.rules = rules
	return nil
}

func (s *AnalyzeDocumentService) GetRulesFilePath() string {
	return s.rulesFile
}

func (s *AnalyzeDocumentService) SetRulesFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("rules file not found: %s", filePath)
	}

	rules, err := models.LoadRulesFromJSON(filePath)
	if err != nil {
		return fmt.Errorf("failed to load rules from file: %w", err)
	}

	s.rulesFile = filePath
	s.rules = rules

	return nil
}

func (s *AnalyzeDocumentService) Execute(text string) *models.ClassificationResult {
	if text == "" {
		return s.createResult("Empty Document", []string{"empty"})
	}

	normalizedText := s.normalizeText(text)

	bestMatchCount := 0
	bestType := ""
	bestKeywords := []string{}

	for _, rule := range s.rules {
		matchedKeywords := []string{}

		for _, keyword := range rule.Keywords {
			normalizedKeyword := strings.ToLower(keyword)
			if strings.Contains(normalizedText, normalizedKeyword) {
				matchedKeywords = append(matchedKeywords, keyword)
			}
		}

		if len(matchedKeywords) > bestMatchCount {
			bestMatchCount = len(matchedKeywords)
			bestType = rule.Type
			bestKeywords = matchedKeywords
		}
	}

	if bestType != "" {
		if len(bestKeywords) > 5 {
			bestKeywords = bestKeywords[:5]
		}

		return s.createResult(bestType, bestKeywords)
	}

	return s.createResult("Other", []string{"document", "text"})
}

func (s *AnalyzeDocumentService) createResult(documentType string, keywords []string) *models.ClassificationResult {
	return &models.ClassificationResult{
		Classification: models.DocumentClassification{
			DocumentType: documentType,
			Keywords:     keywords,
		},
	}
}

func (s *AnalyzeDocumentService) normalizeText(text string) string {
	normalized := strings.ToLower(strings.TrimSpace(text))

	re := regexp.MustCompile(`\s+`)
	normalized = re.ReplaceAllString(normalized, " ")

	return normalized
}
