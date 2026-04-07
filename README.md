# Bow - Copier Parts Cross-Reference

A Go-based tool for parsing Canon copier parts manuals (PDFs) into a SQLite database to enable part cross-referencing between different models.

## Prerequisites

To run this tool, you need the following installed and in your system PATH:

1.  **Go (1.21+)**: [Download and Install Go](https://go.dev/dl/)
2.  **QPDF**: Used for decrypting Canon manuals.
    -   *Install via Winget:* `winget install QPDF.QPDF`
    -   *Install via Choco:* `choco install qpdf`
3.  **PDFtoText**: Part of the Xpdf tools, used for text extraction.
    -   *Windows Tip:* If you have **Git for Windows** installed, it is usually located at `C:\Program Files\Git\mingw64\bin\pdftotext.exe`.
    -   *Note:* The current code uses the Git for Windows path. If yours is different, update the path in `main.go`.

## Setup

1.  Clone the repository.
2.  Create a folder named `Parts/` in the root directory.
3.  Place your Canon Parts Catalog PDF files into the `Parts/` folder.
4.  Run `go mod tidy` to install dependencies.

### 3. Launch the Web Interface
Run the standalone web server to search and cross-reference parts:
```bash
go run server.go
```
Then open your browser to [http://localhost:8080](http://localhost:8080).

## Features
- **Smart Search**: Enter a partial part number (with or without hyphens) or a description.
- **Cross-Reference Grouping**: Results are grouped by Base Part, showing all models that use that part in one view.
- **Detailed Metadata**: See exactly which Figure and Key Number to look for in each model's manual.

## Project Structure
- `main.go`: The core parser and database builder.
- `check_db.go`: A utility script to query the database and show cross-references.
- `parts.db`: The generated SQLite database (ignored by git).
- `Parts/`: Directory for input PDF manuals (ignored by git).
