package models

type ProcessingConfig struct {
	OutputDirectory string
	MoveFiles       bool
}

type ProcessingResult struct {
	ProcessedCount int
	FailedCount    int
	Results        []FileProcessingResult
}

type FileProcessingResult struct {
	Filename     string
	Success      bool
	DocumentType string
	Error        string
}
