package models

// DocumentClassification representa a classificação de um documento
type DocumentClassification struct {
	DocumentType string   `json:"documentType"`
	Keywords     []string `json:"keywords"`
}

// DocumentMetadata contém metadados sobre um documento
type DocumentMetadata struct {
	Filename       string                  `json:"filename"`
	Text           string                  `json:"text"`
	Classification *DocumentClassification `json:"classification,omitempty"`
}

// ClassificationResult representa o resultado da classificação
type ClassificationResult struct {
	Classification DocumentClassification `json:"classification"`
	Summary        string                 `json:"summary,omitempty"`
}
