package services

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"relatorios/interfaces"
	"relatorios/models"
	"relatorios/services/extractors"
)

// DocumentProcessingService gerencia o processamento de documentos
type DocumentProcessingService struct {
	extractorFactory *extractors.DocumentExtractorFactory
	classifier       interfaces.DocumentClassifier
	config           models.ProcessingConfig
}

// NewDocumentProcessingService cria uma nova instância do serviço
func NewDocumentProcessingService(
	extractorFactory *extractors.DocumentExtractorFactory,
	classifier interfaces.DocumentClassifier,
	config models.ProcessingConfig,
) *DocumentProcessingService {
	return &DocumentProcessingService{
		extractorFactory: extractorFactory,
		classifier:       classifier,
		config:           config,
	}
}

// GetClassifierName retorna o nome do classificador
func (s *DocumentProcessingService) GetClassifierName() string {
	return s.classifier.GetClassifierName()
}

// GetSupportedFormats retorna os formatos suportados
func (s *DocumentProcessingService) GetSupportedFormats() []string {
	return s.extractorFactory.GetSupportedFormats()
}

// GetOutputDirectory retorna o diretório de saída configurado
func (s *DocumentProcessingService) GetOutputDirectory() string {
	return s.config.OutputDirectory
}

// ProcessSingleFile processa um único arquivo
func (s *DocumentProcessingService) ProcessSingleFile(filePath string) (models.DocumentMetadata, string, error) {
	// Verificar se o formato é suportado
	if !s.extractorFactory.IsFormatSupported(filePath) {
		return models.DocumentMetadata{}, "", fmt.Errorf("formato não suportado: %s", filepath.Ext(filePath))
	}

	// Obter o extrator adequado para o tipo de arquivo
	extractor, err := s.extractorFactory.GetExtractorForFile(filePath)
	if err != nil {
		return models.DocumentMetadata{}, "", err
	}

	// Extrair texto do documento
	document, err := extractor.ExtractText(filePath)
	if err != nil {
		return models.DocumentMetadata{}, "", fmt.Errorf("falha ao extrair texto: %w", err)
	}

	// Classificar o documento
	document, err = s.classifier.Classify(document)
	if err != nil {
		return models.DocumentMetadata{}, "", fmt.Errorf("falha ao classificar: %w", err)
	}

	// Organizar o arquivo
	destinationPath, err := s.organizeFile(filePath, document.Classification.DocumentType)
	if err != nil {
		return document, "", err
	}

	return document, destinationPath, nil
}

// ProcessDirectory processa todos os arquivos em um diretório
func (s *DocumentProcessingService) ProcessDirectory(dirPath string) (*models.ProcessingResult, error) {
	// Verificar se o diretório existe
	fileInfo, err := os.Stat(dirPath)
	if err != nil {
		return nil, fmt.Errorf("erro ao acessar diretório: %w", err)
	}

	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("o caminho fornecido não é um diretório")
	}

	// Listar todos os arquivos no diretório
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar arquivos no diretório: %w", err)
	}

	// Resultados
	result := &models.ProcessingResult{
		Results: make([]models.FileProcessingResult, 0),
	}

	// Processar cada arquivo
	for _, file := range files {
		if file.IsDir() {
			continue // Ignorar subdiretórios
		}

		filePath := filepath.Join(dirPath, file.Name())

		// Verificar se o formato é suportado
		if !s.extractorFactory.IsFormatSupported(filePath) {
			result.FailedCount++
			result.Results = append(result.Results, models.FileProcessingResult{
				Filename: file.Name(),
				Success:  false,
				Error:    "Formato não suportado",
			})
			continue
		}

		document, _, err := s.ProcessSingleFile(filePath)
		if err != nil {
			result.FailedCount++
			result.Results = append(result.Results, models.FileProcessingResult{
				Filename: file.Name(),
				Success:  false,
				Error:    err.Error(),
			})
		} else {
			result.ProcessedCount++
			result.Results = append(result.Results, models.FileProcessingResult{
				Filename:     file.Name(),
				Success:      true,
				DocumentType: document.Classification.DocumentType,
			})
		}
	}

	return result, nil
}

// organizeFile copia ou move um arquivo para seu diretório de destino com base no tipo
func (s *DocumentProcessingService) organizeFile(filePath string, documentType string) (string, error) {
	// Criar o diretório de destino se não existir
	destDir := filepath.Join(s.config.OutputDirectory, documentType)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("erro ao criar diretório de destino: %w", err)
	}

	// Nome do arquivo de destino
	fileName := filepath.Base(filePath)
	destPath := filepath.Join(destDir, fileName)

	// Copiar ou mover o arquivo
	if s.config.MoveFiles {
		if err := os.Rename(filePath, destPath); err != nil {
			return "", fmt.Errorf("erro ao mover arquivo: %w", err)
		}
	} else {
		srcFile, err := os.Open(filePath)
		if err != nil {
			return "", fmt.Errorf("erro ao abrir arquivo de origem: %w", err)
		}
		defer srcFile.Close()

		destFile, err := os.Create(destPath)
		if err != nil {
			return "", fmt.Errorf("erro ao criar arquivo de destino: %w", err)
		}
		defer destFile.Close()

		if _, err := io.Copy(destFile, srcFile); err != nil {
			return "", fmt.Errorf("erro ao copiar arquivo: %w", err)
		}
	}

	return destPath, nil
}
