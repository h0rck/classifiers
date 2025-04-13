package services

import (
	"regexp"
	"relatorios/models"
	"strings"
)

// AnalyzeDocumentService é responsável por analisar e classificar documentos
type AnalyzeDocumentService struct {
	// Adicionar quaisquer dependências necessárias
}

// NewAnalyzeDocumentService cria uma nova instância do serviço
func NewAnalyzeDocumentService() *AnalyzeDocumentService {
	return &AnalyzeDocumentService{}
}

// Execute analisa o texto do documento e gera uma classificação
func (s *AnalyzeDocumentService) Execute(text string) *models.ClassificationResult {
	// Normalizar o texto para comparação
	normalizedText := s.normalizeText(text)

	// Verificar tipos de documentos baseados em padrões
	if s.matchesContractPattern(normalizedText) {
		return s.createResult("Contrato", 0.92, []string{
			"contrato", "partes", "cláusula", "rescisão", "obrigações"})
	} else if s.matchesInvoicePattern(normalizedText) {
		return s.createResult("Nota Fiscal", 0.95, []string{
			"nota fiscal", "cnpj", "valor total", "emissão", "imposto"})
	} else if s.matchesReceiptPattern(normalizedText) {
		return s.createResult("Recibo", 0.90, []string{
			"recibo", "pagamento", "valor", "recebemos", "quitar"})
	} else if s.matchesReportPattern(normalizedText) {
		return s.createResult("Relatório", 0.85, []string{
			"relatório", "análise", "conclusão", "resultados", "período"})
	} else if s.matchesCurriculumPattern(normalizedText) {
		return s.createResult("Currículo", 0.89, []string{
			"currículo", "experiência", "formação", "habilidades", "profissional"})
	}

	// Documento não reconhecido
	return s.createResult("Outro", 0.50, []string{"documento", "texto"})
}

// Funções auxiliares para reconhecer tipos de documentos
func (s *AnalyzeDocumentService) matchesContractPattern(text string) bool {
	contractKeywords := []string{"contrato", "cláusula", "partes", "rescisão", "acordo", "contratante", "contratado"}
	return s.containsMultipleKeywords(text, contractKeywords, 3)
}

func (s *AnalyzeDocumentService) matchesInvoicePattern(text string) bool {
	// Verificar padrões de nota fiscal
	invoiceKeywords := []string{"nota fiscal", "nf-e", "cnpj", "discriminação", "valor total", "data de emissão"}
	hasKeywords := s.containsMultipleKeywords(text, invoiceKeywords, 3)

	// Verificar padrão de CNPJ ou CPF
	cnpjPattern := regexp.MustCompile(`\d{2}\.\d{3}\.\d{3}\/\d{4}\-\d{2}`)
	cpfPattern := regexp.MustCompile(`\d{3}\.\d{3}\.\d{3}\-\d{2}`)
	hasDocNumber := cnpjPattern.MatchString(text) || cpfPattern.MatchString(text)

	return hasKeywords && hasDocNumber
}

func (s *AnalyzeDocumentService) matchesReceiptPattern(text string) bool {
	receiptKeywords := []string{"recibo", "recebi", "valor", "quantia", "pagamento", "quita"}
	return s.containsMultipleKeywords(text, receiptKeywords, 2)
}

func (s *AnalyzeDocumentService) matchesReportPattern(text string) bool {
	reportKeywords := []string{"relatório", "análise", "conclusão", "avaliação", "resultados", "período"}
	return s.containsMultipleKeywords(text, reportKeywords, 2)
}

func (s *AnalyzeDocumentService) matchesCurriculumPattern(text string) bool {
	cvKeywords := []string{"currículo", "curriculum", "vitae", "experiência", "formação", "profissional", "habilidades"}
	return s.containsMultipleKeywords(text, cvKeywords, 3)
}

// Funções utilitárias

// normalizeText normaliza o texto para comparação
func (s *AnalyzeDocumentService) normalizeText(text string) string {
	text = strings.ToLower(text)
	text = s.removeAccents(text)
	return text
}

// removeAccents remove acentos de um texto
func (s *AnalyzeDocumentService) removeAccents(text string) string {
	replacements := map[string]string{
		"á": "a", "à": "a", "â": "a", "ã": "a", "ä": "a",
		"é": "e", "è": "e", "ê": "e", "ë": "e",
		"í": "i", "ì": "i", "î": "i", "ï": "i",
		"ó": "o", "ò": "o", "ô": "o", "õ": "o", "ö": "o",
		"ú": "u", "ù": "u", "û": "u", "ü": "u",
		"ç": "c", "ñ": "n",
	}

	for accentedChar, plainChar := range replacements {
		text = strings.ReplaceAll(text, accentedChar, plainChar)
	}
	return text
}

// containsMultipleKeywords verifica se o texto contém múltiplas palavras-chave
func (s *AnalyzeDocumentService) containsMultipleKeywords(text string, keywords []string, minMatches int) bool {
	matches := 0
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			matches++
			if matches >= minMatches {
				return true
			}
		}
	}
	return false
}

// createResult cria um objeto de resultado de classificação
func (s *AnalyzeDocumentService) createResult(documentType string, confidence float64, keywords []string) *models.ClassificationResult {
	return &models.ClassificationResult{
		Classification: models.DocumentClassification{
			DocumentType: documentType,
			Confidence:   confidence,
			Keywords:     keywords,
		},
	}
}
