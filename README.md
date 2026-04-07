# CanonBow - Copier Parts Cross-Reference System

CanonBow is a high-performance Go-based utility designed to ingest Canon Parts Catalog PDFs and build a searchable SQLite database. It specializes in "Smart Cross-Referencing," allowing technicians to see exactly which parts are shared across different copier series (imageRUNNER ADVANCE, imageFORCE, imagePRESS, etc.).

## 🚀 Project Achievements & Features
- **High-Speed Parsing**: Optimized extraction using `pdftotext -table` mode, processing 70+ manuals in under 10 minutes.
- **Smart Metadata Extraction**: Automatically detects Model Series and Catalog Revision from both filenames and PDF content.
- **Robust Regex Logic**: Cleanly separates Key Numbers, Base Part Numbers, Revisions, and Descriptions, even when PDF text is tightly packed.
- **Unified Cross-Referencing**: Databases are structured by "Base Part Number" (first 8 chars), ignoring revision numbers to find compatibility across models.
- **Web Interface**: A standalone `localhost:8080` server with grouping logic to show all model occurrences for any searched part.
- **Safe Concurrency**: The web server uses SQLite Read-Only (RO) mode with WAL journaling to prevent database locks while the builder is running.

---

## 🛠 Prerequisites & Tools
To run or develop this project on a new machine, ensure the following are installed:

### 1. Go (Golang)
- **Version**: 1.21 or higher.
- **Install**: [go.dev/dl](https://go.dev/dl/)

### 2. QPDF (Encryption Handler)
- **Purpose**: Canon manuals are often encrypted; QPDF decrypts them to a temporary state for parsing.
- **Install**: `winget install QPDF.QPDF`
- **Verification**: Ensure `qpdf.exe` is in your PATH.

### 3. PDFtoText (Text Extractor)
- **Purpose**: High-fidelity text extraction using `-table` mode to preserve column alignment.
- **Source**: Part of Git for Windows (MinGW64).
- **Default Path**: `C:\Program Files\Git\mingw64\bin\pdftotext.exe`
- **Manual Path Update**: If installed elsewhere, update the path in `main.go` and `server.go`.

### 4. SQLite3 CLI (Optional for Debugging)
- **Purpose**: Direct querying of `parts.db`.
- **Install**: `winget install SQLite.SQLite`

---

## 📂 Project Structure
- `main.go`: The Database Builder. Scans `Parts/`, decrypts, and parses content.
- `server.go`: The Web Front-end. Provides the UI and smart search logic.
- `check_db.go`: CLI utility for quick database stats and cross-ref samples.
- `investigate_db.go`: Tool for finding "junk" data and refining regex.
- `Parts/`: Drop your `.pdf` manuals here (git-ignored, but the folder is tracked).
- `parts.db`: The local SQLite database (git-ignored).

---

## 📖 Instructions for You (Jeremy)

### How to add new manuals:
1.  Drop the new PDFs into the `Parts/` folder.
2.  Run the builder: `go run main.go`.
3.  The tool will automatically detect the series names and add the parts to your database.

### How to use the Web UI:
1.  Run the server: `go run server.go`.
2.  Open [http://localhost:8080](http://localhost:8080).
3.  Search for a part number (e.g., `WG8-5935` or `FE31525`) or a keyword (e.g., `Roller`).

---

## 📝 Developer Notes (For Gemini CLI)
- **Regex Strategy**: The primary matcher is `(?m)(?:^|\s{2,})(\d{1,3})\s{2,}([A-Z0-9]{3}-[A-Z0-9]{4})-([A-Z0-9]{3})\s{1,}(\d*)\s+(.*?)(\s{2,}|\r|\n|$)`.
- **Database Schema**: Uses `manuals`, `figures`, and `parts` tables. `base_part` is indexed for $O(1)$ search performance.
- **Filters**: Explicitly ignores `SNL` (Serial Number) figures and descriptions < 3 characters to maintain data purity.
