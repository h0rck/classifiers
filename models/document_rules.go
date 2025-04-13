package models

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// DocumentRule define uma regra para classificação de documentos
type DocumentRule struct {
	Type     string   `json:"type"`
	Keywords []string `json:"keywords"`
}

// LoadRulesFromJSON carrega regras de classificação de um arquivo JSON
func LoadRulesFromJSON(filePath string) ([]DocumentRule, error) {
	// Verificar se o arquivo existe
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Se o arquivo não existe, criar com as regras padrão
		defaultRules := GetDefaultRules()
		if err := SaveRulesToJSON(filePath, defaultRules); err != nil {
			return nil, fmt.Errorf("falha ao criar arquivo de regras padrão: %w", err)
		}
		return defaultRules, nil
	}

	// Ler o arquivo JSON
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler arquivo de regras: %w", err)
	}

	// Decodificar o JSON
	var rules []DocumentRule
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, fmt.Errorf("falha ao decodificar regras JSON: %w", err)
	}

	return rules, nil
}

// SaveRulesToJSON salva as regras em um arquivo JSON
func SaveRulesToJSON(filePath string, rules []DocumentRule) error {
	// Garantir que o diretório existe
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("falha ao criar diretório para o arquivo de regras: %w", err)
	}

	// Codificar as regras para JSON com formatação legível
	data, err := json.MarshalIndent(rules, "", "  ")
	if err != nil {
		return fmt.Errorf("falha ao codificar regras para JSON: %w", err)
	}

	// Escrever no arquivo
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("falha ao salvar arquivo de regras: %w", err)
	}

	return nil
}

// GetDefaultRules retorna o conjunto padrão de regras de classificação
func GetDefaultRules() []DocumentRule {
	return []DocumentRule{
		{
			Type: "Contrato",
			Keywords: []string{
				"contrato", "cláusula", "partes", "rescisão", "acordo",
				"contratante", "contratado", "obrigações", "vigência",
				"objeto", "firmado", "assinatura", "foro", "jurisdição",
				"prazo", "pagamento", "multas", "condições", "confidencialidade",
			},
		},
		{
			Type: "Nota Fiscal",
			Keywords: []string{
				"nota fiscal", "nf-e", "nfe", "cnpj", "emissão",
				"impostos", "valor total", "data de emissão", "discriminação",
				"produto", "quantidade", "total", "icms", "ipi",
				"cofins", "pis", "alíquota", "natureza da operação",
				"destinatário", "emitente",
			},
		},
		{
			Type: "Recibo",
			Keywords: []string{
				"recibo", "recebi", "valor", "quantia", "pagamento",
				"referente", "importância", "pago", "assinatura",
				"recebedor", "pagador", "comprovante", "quitado",
				"data do pagamento",
			},
		},
		{
			Type: "Relatório",
			Keywords: []string{
				"relatório", "análise", "conclusão", "avaliação",
				"resultados", "período", "dados", "pesquisa",
				"metodologia", "introdução", "objetivo", "sumário",
				"estatísticas", "gráficos", "observações",
			},
		},
		{
			Type: "Currículo",
			Keywords: []string{
				"currículo", "curriculum", "vitae", "experiência",
				"formação", "profissional", "habilidades", "escolaridade",
				"idiomas", "qualificações", "certificações",
				"conhecimentos", "objetivo profissional", "referências",
				"contato", "telefone", "email", "linkedin",
			},
		},
	}
}
