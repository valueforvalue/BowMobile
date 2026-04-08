# CanonBow - Copier Parts Cross-Reference System

CanonBow is a high-performance Go-based utility designed to ingest Canon Parts Catalog PDFs and build a searchable SQLite database. It specializes in "Smart Cross-Referencing," allowing technicians to see exactly which parts are shared across different copier series (imageRUNNER ADVANCE, imageFORCE, imagePRESS, etc.).

## 🚀 Project Achievements & Features
- **High-Speed Parsing**: Optimized extraction using a multi-step pipeline (Noise Filter -> Anchor Search -> Field Extraction -> Validation).
- **Modern Web Stack**: Migrated to **Go + Templ + HTMX + Tailwind CSS** for a fast, "live" search experience.
- **Smart Remarks Extraction**: Automatically captures schematic locations and connectors (e.g., **SL1, PS64, J705**) from manuals.
- **Print Selection**: Select specific part occurrences and generate a clean, formatted parts list for printing.
- **Incremental Builds**: The builder automatically skips manuals already in the database, making updates near-instant.
- **Portable Release**: Single-file **`Bow.exe`** with embedded database, auto-browser launch, and auto-shutdown heartbeat.

---

## 🛠 Prerequisites & Tools
To develop this project or rebuild from source, ensure the following are installed:

### 1. Go (Golang) & Templ
- **Go**: v1.21 or higher. [go.dev/dl](https://go.dev/dl/)
- **Templ**: `go install github.com/a-h/templ/cmd/templ@latest`

### 2. QPDF (Encryption Handler)
- **Install**: `winget install QPDF.QPDF`
- **Verification**: Ensure `qpdf.exe` is in your PATH.

### 3. PDFtoText (Text Extractor)
- **Source**: Part of Git for Windows (MinGW64).
- **Default Path**: `C:\Program Files\Git\mingw64\bin\pdftotext.exe`

---

## 📂 Project Structure
- `models.go`: Shared types and constants (the `bow` package).
- `cmd/server/`: The Web UI server (Templ + HTMX).
- `cmd/builder/`: The PDF parsing and database building engine.
- `cmd/tools/`: Diagnostic and utility scripts (e.g., `check_db.go`).
- `Parts/`: Directory for source PDF manuals (git-ignored).
- `assets/`: UI assets like the project logo.
- `build_release.ps1`: Automation script for creating the standalone ZIP release.

---

## 📖 Instructions

### How to develop/run from source:
1.  **Generate Templates**: `templ generate ./cmd/server/`
2.  **Start Server**: `go run ./cmd/server`
3.  **Run Builder**: `go run ./cmd/builder` (add a PDF path as an argument for single-file mode).

### How to create a new release:
1.  Ensure `parts.db` is populated.
2.  Run `./build_release.ps1`.
3.  This creates `Bow.exe` and `Bow_Release.zip` in the root folder.

---

## 🗺️ Roadmap
- [ ] **Accessory Expansion**: Ingest part catalogs for finishers, paper decks, and other peripherals.
- [ ] **Cloud Sync**: Option for the client to check GitHub for an updated `parts.db` file automatically.
- [ ] **Image Integration**: Explore linking part numbers to extracted diagram images from the PDFs.

---

## 📝 Developer Notes (For Gemini CLI)
- **Heartbeat System**: The server uses a `/pulse` endpoint. If the browser tab is closed for > 30s, the process exits automatically to free up port 8080.
- **Embedded Files**: The release build embeds `parts.db` and `assets/`. On startup, `Bow.exe` extracts the DB to the current directory if it doesn't exist.
- **Search Logic**: Supports both Part Number (normalized) and keyword searches across both `description` and `remarks` columns.
