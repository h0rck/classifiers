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

**Windows**: Download from [UB-Mannheim/tesseract](https://github.com/UB-Mannheim/tesseract/wiki)

**macOS**: `brew install tesseract tesseract-lang`

**Linux (Ubuntu/Debian)**:
```
sudo apt-get update && sudo apt-get install -y tesseract-ocr tesseract-ocr-eng
```

## Usage
Run with path: `./classifiers`




