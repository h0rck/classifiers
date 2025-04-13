package services

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"relatorios/interfaces"
	"relatorios/models"
	"relatorios/services/classifiers"
	"relatorios/services/extractors"
)

type DocumentProcessingService struct {
	extractorFactory *extractors.DocumentExtractorFactory
	classifier       interfaces.DocumentClassifier
	config           models.ProcessingConfig
}

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

func (s *DocumentProcessingService) GetClassifierName() string {
	return s.classifier.GetClassifierName()
}

func (s *DocumentProcessingService) GetSupportedFormats() []string {
	return s.extractorFactory.GetSupportedFormats()
}

func (s *DocumentProcessingService) GetOutputDirectory() string {
	return s.config.OutputDirectory
}

func (s *DocumentProcessingService) ProcessSingleFile(filePath string) (models.DocumentMetadata, string, error) {
	if !s.extractorFactory.IsFormatSupported(filePath) {
		return models.DocumentMetadata{}, "", fmt.Errorf("unsupported format: %s", filepath.Ext(filePath))
	}

	extractor, err := s.extractorFactory.GetExtractorForFile(filePath)
	if err != nil {
		return models.DocumentMetadata{}, "", err
	}

	document, err := extractor.ExtractText(filePath)
	if err != nil {
		return models.DocumentMetadata{}, "", fmt.Errorf("failed to extract text: %w", err)
	}

	document, err = s.classifier.Classify(document)
	if err != nil {
		return models.DocumentMetadata{}, "", fmt.Errorf("classification failed: %w", err)
	}

	destinationPath, err := s.organizeFile(filePath, document.Classification.DocumentType)
	if err != nil {
		return document, "", err
	}

	return document, destinationPath, nil
}

func (s *DocumentProcessingService) ProcessDirectory(dirPath string) (*models.ProcessingResult, error) {
	fileInfo, err := os.Stat(dirPath)
	if err != nil {
		return nil, fmt.Errorf("error accessing directory: %w", err)
	}

	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("provided path is not a directory")
	}

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("error listing files in directory: %w", err)
	}

	result := &models.ProcessingResult{
		Results: make([]models.FileProcessingResult, 0),
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(dirPath, file.Name())

		if !s.extractorFactory.IsFormatSupported(filePath) {
			result.FailedCount++
			result.Results = append(result.Results, models.FileProcessingResult{
				Filename: file.Name(),
				Success:  false,
				Error:    "Unsupported format",
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

func (s *DocumentProcessingService) organizeFile(filePath string, documentType string) (string, error) {
	destDir := filepath.Join(s.config.OutputDirectory, documentType)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("error creating destination directory: %w", err)
	}

	fileName := filepath.Base(filePath)
	destPath := filepath.Join(destDir, fileName)

	if s.config.MoveFiles {
		if err := os.Rename(filePath, destPath); err != nil {
			return "", fmt.Errorf("error moving file: %w", err)
		}
	} else {
		srcFile, err := os.Open(filePath)
		if err != nil {
			return "", fmt.Errorf("error opening source file: %w", err)
		}
		defer srcFile.Close()

		destFile, err := os.Create(destPath)
		if err != nil {
			return "", fmt.Errorf("error creating destination file: %w", err)
		}
		defer destFile.Close()

		if _, err := io.Copy(destFile, srcFile); err != nil {
			return "", fmt.Errorf("error copying file: %w", err)
		}
	}

	return destPath, nil
}

func (s *DocumentProcessingService) GetAnalyzeService() interfaces.AnalyzeService {
	if classifier, ok := s.classifier.(*classifiers.DocumentClassifier); ok {
		return classifier.GetAnalyzeService()
	}
	return nil
}
