package models

type DocumentClassification struct {
	DocumentType string   `json:"documentType"`
	Keywords     []string `json:"keywords"`
}

type DocumentMetadata struct {
	Filename       string                  `json:"filename"`
	Text           string                  `json:"text"`
	Classification *DocumentClassification `json:"classification,omitempty"`
}

type ClassificationResult struct {
	Classification DocumentClassification `json:"classification"`
	Summary        string                 `json:"summary,omitempty"`
}
