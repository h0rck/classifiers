package models

type DocumentRule struct {
	Type     string
	Keywords []string
}

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
