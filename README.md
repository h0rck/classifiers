# classifiers

A powerful document processing system written in Go that extracts data from documents and classifies them according to configurable rules.

## Overview

classifiers is designed to streamline document processing workflows by automating the extraction of relevant information from various document types and classifying them based on customizable rule sets. It features a console interface for easy interaction and uses a configurable architecture to adapt to different document processing needs.

## ✨ How It Works ✨

The document classification system automatically categorizes documents based on their content and predefined rules:

### 1. Classification Rules
> Rules are defined with a document type and a set of keywords that identify that type.
> The system matches these keywords against document content to determine its category.

### 2. Document Processing
> **Extract** → **Classify** → **Organize**
> 
> - **Extract**: The system pulls text content from various document formats
> - **Classify**: It applies classification rules to determine document type
> - **Organize**: Documents are sorted according to their classification

### 3. File Browser Interface
> Navigate your file system with ease:
> 
> - 📁 Browse directories with a clear visual representation
> - 📄 Select individual files for processing
> - 🔍 Filter by supported file types
> - ⌨️ Enter paths manually when needed

### 4. Classification Rules Management
> Take control of your document rules:
> 
> - 📋 View and edit current classification rules
> - 🔄 Reload rules from external files
> - 🔀 Select different rule sets for different document types

### Supported Document Types

The classifier can identify various document types including:
- 📊 Invoices
- 📜 Contracts
- 🧾 Receipts
- 📝 Reports
- ➕ And more, depending on your rule configurations



## Project Structure

- `/models` - Data structure definitions for document processing
- `/services` - Core business logic implementation
- `/services/extractors` - Document data extraction components
- `/services/classifiers` - Document classification logic
- `/ui` - User interface components

## Installation

### Prerequisites

- Go 1.16 or higher
- Git

### Steps

1. Clone the repository:
   ```
   git clone https://github.com/h0rck/classifiers-go.git
   cd classifiers
   ```

2. Build the application:
   ```
   go build
   ```

## Dependencies

### Tesseract OCR

This application requires Tesseract OCR for image text extraction. Follow these instructions to install it:

#### Windows
1. Visit https://github.com/UB-Mannheim/tesseract/wiki
2. Download the latest installer (tesseract-ocr-w64-setup-v5.x.x.exe)
3. During installation, select desired languages
4. Restart the application after installation

#### macOS
1. Install Homebrew if you don't have it (https://brew.sh/)
2. Run in terminal: `brew install tesseract tesseract-lang`
3. Restart the application after installation

#### Linux
Ubuntu/Debian:
```
sudo apt-get update && sudo apt-get install -y tesseract-ocr tesseract-ocr-eng
```

Fedora/RHEL:
```
sudo dnf install -y tesseract tesseract-langpack-eng
```

## Usage

### Basic Usage

Run the application with an optional path to process:

```
./classifiers
```

If no path is provided, the application will prompt for one.

### Command Line Arguments

- First argument (optional): Path to the directory containing documents to process

## Configuration

The application stores its configuration in the user's configuration directory:

- Linux: `~/.config/classifiers/`
- macOS: `~/Library/Application Support/classifiers/`
- Windows: `%AppData%\classifiers\`

### Document Rules

Document classification rules are stored in `document_rules.json`. This file is automatically created in the configuration directory when the application is first run.

## Customization

### Processing Configuration

You can modify the default processing behavior by adjusting the `models.ProcessingConfig` in the main.go file:

- `OutputDirectory`: Directory where processed documents will be saved
- `MoveFiles`: Whether to move original files after processing

### Adding New Extractors

To support additional document types, create new extractors in the `/services/extractors` directory and register them in the `DocumentExtractorFactory`.



