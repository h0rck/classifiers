package models

// ProcessingConfig contém configurações para o processamento de documentos
type ProcessingConfig struct {
	OutputDirectory string
	MoveFiles       bool
}

// ProcessingResult representa os resultados do processamento de vários documentos
type ProcessingResult struct {
	ProcessedCount int
	FailedCount    int
	Results        []FileProcessingResult
}

// FileProcessingResult representa o resultado do processamento de um arquivo
type FileProcessingResult struct {
	Filename     string
	Success      bool
	DocumentType string
	Error        string
}
