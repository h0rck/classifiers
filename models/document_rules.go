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
		// Se o arquivo não existe, exibir instruções para criá-lo
		fmt.Printf("\n===============================================================\n")
		fmt.Printf("AVISO: Arquivo de regras não encontrado: %s\n\n", filePath)
		fmt.Printf("Para usar o classificador, crie um arquivo JSON com o seguinte formato:\n\n")
		fmt.Printf("[\n")
		fmt.Printf("  {\n")
		fmt.Printf("    \"type\": \"Nome do Tipo de Documento\",\n")
		fmt.Printf("    \"keywords\": [\n")
		fmt.Printf("      \"palavra-chave1\",\n")
		fmt.Printf("      \"palavra-chave2\",\n")
		fmt.Printf("      \"frase chave também funciona\"\n")
		fmt.Printf("    ]\n")
		fmt.Printf("  },\n")
		fmt.Printf("  {\n")
		fmt.Printf("    \"type\": \"Outro Tipo de Documento\",\n")
		fmt.Printf("    \"keywords\": [\n")
		fmt.Printf("      \"outra-palavra-chave\",\n")
		fmt.Printf("      \"mais uma palavra\"\n")
		fmt.Printf("    ]\n")
		fmt.Printf("  }\n")
		fmt.Printf("]\n\n")
		fmt.Printf("Salve este arquivo em: %s\n", filePath)
		fmt.Printf("===============================================================\n\n")

		// Exemplo prático para ajudar o usuário
		fmt.Printf("Exemplo para classificar documentos comuns:\n\n")
		fmt.Printf("[\n")
		fmt.Printf("  {\n")
		fmt.Printf("    \"type\": \"Contrato\",\n")
		fmt.Printf("    \"keywords\": [\n")
		fmt.Printf("      \"contrato\", \"cláusula\", \"partes\", \"rescisão\",\n")
		fmt.Printf("      \"acordo\", \"contratante\", \"vigência\"\n")
		fmt.Printf("    ]\n")
		fmt.Printf("  },\n")
		fmt.Printf("  {\n")
		fmt.Printf("    \"type\": \"Nota Fiscal\",\n")
		fmt.Printf("    \"keywords\": [\n")
		fmt.Printf("      \"nota fiscal\", \"nf-e\", \"cnpj\", \"emissão\",\n")
		fmt.Printf("      \"valor total\", \"icms\", \"data de emissão\"\n")
		fmt.Printf("    ]\n")
		fmt.Printf("  },\n")
		fmt.Printf("  {\n")
		fmt.Printf("    \"type\": \"Recibo\",\n")
		fmt.Printf("    \"keywords\": [\n")
		fmt.Printf("      \"recibo\", \"recebi\", \"valor\", \"quantia\",\n")
		fmt.Printf("      \"pagamento\", \"importância\", \"pago\"\n")
		fmt.Printf("    ]\n")
		fmt.Printf("  }\n")
		fmt.Printf("]\n\n")

		// Criar uma regra simples para começar
		simpleRules := []DocumentRule{
			{
				Type:     "Documento",
				Keywords: []string{"texto", "documento"},
			},
		}

		fmt.Printf("Criando um arquivo básico para você começar...\n")

		// Garantir que o diretório existe
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return simpleRules, nil
		}

		// Salvar um arquivo básico
		if err := SaveRulesToJSON(filePath, simpleRules); err != nil {
			fmt.Printf("Não foi possível criar o arquivo de regras. Usando regras mínimas.\n")
			return simpleRules, nil
		}

		fmt.Printf("Arquivo de regras básico criado em: %s\n", filePath)
		fmt.Printf("Edite este arquivo para personalizar sua classificação.\n")
		fmt.Printf("===============================================================\n")

		return simpleRules, nil
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
