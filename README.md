# classifiers

A powerful document processing system written in Go that extracts data from documents and classifies them according to configurable rules.

## Overview

classifiers is designed to streamline document processing workflows by automating the extraction of relevant information from various document types and classifying them based on customizable rule sets. It features a console interface for easy interaction and uses a configurable architecture to adapt to different document processing needs.

## Features

- Automated document data extraction
- Rule-based document classification
- User-friendly console interface
- Configurable processing workflow
- Document archiving and organization

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

