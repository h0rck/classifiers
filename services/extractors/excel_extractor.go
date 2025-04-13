package extractors

import (
	"fmt"
	"path/filepath"
	"relatorios/models"
	"strings"

	"github.com/xuri/excelize/v2"
)

// ExcelExtractor é um extrator que lê arquivos Excel
type ExcelExtractor struct{}

// ExtractText extrai o texto de arquivos Excel
func (e *ExcelExtractor) ExtractText(filePath string) (models.DocumentMetadata, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return models.DocumentMetadata{}, fmt.Errorf("falha ao abrir arquivo Excel: %w", err)
	}
	defer f.Close()

	var textContent strings.Builder

	for _, sheetName := range f.GetSheetList() {
		textContent.WriteString(fmt.Sprintf("\n[Planilha: %s]\n", sheetName))

		rows, err := f.GetRows(sheetName)
		if err != nil {
			return models.DocumentMetadata{}, fmt.Errorf("falha ao ler planilha %s: %w", sheetName, err)
		}

		for _, row := range rows {
			if len(row) > 0 {
				textContent.WriteString(strings.Join(row, " | ") + "\n")
			}
		}
	}

	return models.DocumentMetadata{
		Filename: filepath.Base(filePath),
		Text:     textContent.String(),
	}, nil
}

// IsSupportedFormat verifica se o formato do arquivo é suportado
func (e *ExcelExtractor) IsSupportedFormat(filePath string) bool {
	ext := filepath.Ext(filePath)
	return ext == ".xlsx" || ext == ".xls"
}

// GetSupportedFormats retorna os formatos suportados por este extrator
func (e *ExcelExtractor) GetSupportedFormats() []string {
	return []string{".xlsx", ".xls"}
}
