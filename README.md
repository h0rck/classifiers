# classifiers

A document processing system written in Go that extracts data from documents and classifies them based on configurable rules.

## How It Works 

### 1. Classification Rules
> Rules define document types and the keywords that identify them.

### 2. Document Processing
> **Extract** → **Classify** → **Organize**
> 
> - **Extract**: Pull text from various document formats
> - **Classify**: Apply rules to determine document type
> - **Organize**: Sort documents by classification

### 3. File Browser Interface
> - 📁 Browse directories with visual representation
> - 📄 Select files for processing
> - 🔍 Filter by supported types

### 4. Classification Rules Management
> - 📋 View and edit rules
> - 🔄 Reload from external files
> - 🔀 Select different rule sets

### Supported Document Types
- 📊 Invoices
- 📜 Contracts
- 🧾 Receipts
- 📝 Reports

## Project Structure
- `/models` - Data structure definitions
- `/services` - Core business logic
- `/services/extractors` - Document data extraction components
- `/services/classifiers` - Document classification logic
- `/ui` - User interface components

## Dependencies

### Tesseract OCR
Required for image text extraction:

**Windows**: 
- Download from [UB-Mannheim/tesseract](https://github.com/UB-Mannheim/tesseract/wiki)
- During installation, select the languages you need (English is 'eng', Portuguese is 'por')

**macOS**: 
```
brew install tesseract tesseract-lang
```
- For specific languages: `brew install tesseract-lang-eng tesseract-lang-por`

**Linux (Ubuntu/Debian)**:
```
# For English language support
sudo apt update && sudo apt install -y tesseract-ocr tesseract-ocr-eng

# For Portuguese language support
sudo apt install -y tesseract-ocr-por
```

**Language Packs**:
- 'eng': English
- 'por': Portuguese
- 'spa': Spanish
- 'fra': French
- 'deu': German

Install the language packs appropriate for your documents' content.

## Usage
Run with path: `./classifiers`




