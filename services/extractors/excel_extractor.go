package extractors

import (
	"fmt"
	"path/filepath"
	"relatorios/models"
	"strings"

	"github.com/xuri/excelize/v2"
)

type ExcelExtractor struct{}

func (e *ExcelExtractor) ExtractText(filePath string) (models.DocumentMetadata, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return models.DocumentMetadata{}, fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer f.Close()

	var textContent strings.Builder

	for _, sheetName := range f.GetSheetList() {
		textContent.WriteString(fmt.Sprintf("\n[Sheet: %s]\n", sheetName))

		rows, err := f.GetRows(sheetName)
		if err != nil {
			return models.DocumentMetadata{}, fmt.Errorf("failed to read sheet %s: %w", sheetName, err)
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

func (e *ExcelExtractor) IsSupportedFormat(filePath string) bool {
	ext := filepath.Ext(filePath)
	return ext == ".xlsx" || ext == ".xls"
}

func (e *ExcelExtractor) GetSupportedFormats() []string {
	return []string{".xlsx", ".xls"}
}
