package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"relatorios/models"
	"relatorios/services"
	"relatorios/services/classifiers"
	"relatorios/services/extractors"
	"relatorios/ui"
)

func main() {
	// Definir o caminho para o arquivo de regras
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		userConfigDir = "."
	}

	configDir := filepath.Join(userConfigDir, "relatorios-go")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Aviso: Não foi possível criar diretório de configuração: %v\n", err)
		configDir = "."
	}

	rulesFile := filepath.Join(configDir, "document_rules.json")

	// Criar regras padrão
	defaultRules := []map[string]interface{}{
		{
			"type": "Contrato",
			"keywords": []string{
				"contrato", "cláusula", "partes", "rescisão", "acordo",
				"contratante", "contratado", "obrigações", "vigência",
				"objeto", "firmado", "assinatura", "foro", "jurisdição",
				"prazo", "pagamento", "multas", "condições", "confidencialidade",
			},
		},
		{
			"type": "Nota Fiscal",
			"keywords": []string{
				"nota fiscal", "nf-e", "nfe", "cnpj", "emissão",
				"impostos", "valor total", "data de emissão", "discriminação",
				"produto", "quantidade", "total", "icms", "ipi",
				"cofins", "pis", "alíquota", "natureza da operação",
				"destinatário", "emitente",
			},
		},
		{
			"type": "Recibo",
			"keywords": []string{
				"recibo", "recebi", "valor", "quantia", "pagamento",
				"referente", "importância", "pago", "assinatura",
				"recebedor", "pagador", "comprovante", "quitado",
				"data do pagamento",
			},
		},
		{
			"type": "Relatório",
			"keywords": []string{
				"relatório", "análise", "conclusão", "avaliação",
				"resultados", "período", "dados", "pesquisa",
				"metodologia", "introdução", "objetivo", "sumário",
				"estatísticas", "gráficos", "observações",
			},
		},
		{
			"type": "Currículo",
			"keywords": []string{
				"currículo", "curriculum", "vitae", "experiência",
				"formação", "profissional", "habilidades", "escolaridade",
				"idiomas", "qualificações", "certificações",
				"conhecimentos", "objetivo profissional", "referências",
				"contato", "telefone", "email", "linkedin",
			},
		},
	}

	// Salvar regras padrão no arquivo de regras, se ele não existir
	if _, err := os.Stat(rulesFile); os.IsNotExist(err) {
		file, err := os.Create(rulesFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao criar arquivo de regras: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(defaultRules); err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao salvar regras padrão: %v\n", err)
			os.Exit(1)
		}
	}

	// Inicializar componentes do sistema
	extractorFactory := extractors.NewDocumentExtractorFactory()
	analyzeDocumentService := services.NewAnalyzeDocumentService(rulesFile)
	classifier := classifiers.NewDocumentClassifier(analyzeDocumentService)

	// Configurar o serviço de processamento de documentos
	config := models.ProcessingConfig{
		OutputDirectory: "./output",
		MoveFiles:       false,
	}

	processingService := services.NewDocumentProcessingService(
		extractorFactory,
		classifier,
		config,
	)

	// Criar a interface de console
	consoleInterface := ui.NewConsoleInterface(processingService)

	// Determinar se o programa foi chamado com um argumento de caminho
	var initialPath string
	if len(os.Args) > 1 {
		initialPath = os.Args[1]
	}

	// Iniciar a interface
	err = consoleInterface.Start(initialPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao iniciar a aplicação: %v\n", err)
		os.Exit(1)
	}
}
