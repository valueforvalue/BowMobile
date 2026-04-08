# Bow - Copier Parts Cross-Reference

## Project Overview
Bow is a desktop application (built with Go and Wails) used for cross-referencing copier parts across different manuals and models. It uses a SQLite database (`parts.db`) stored as a sidecar file.

## Tech Stack
- **Backend**: Go
- **GUI Framework**: [Wails](https://wails.io/) (v2)
- **Database**: SQLite (using `modernc.org/sqlite`)
- **Templates**: [Templ](https://templ.guide/)
- **Frontend**: Vite, Tailwind CSS, Vanilla JS

## Critical Files
- `bow-gui/app.go`: Main application logic and backend bindings for Wails.
- `cmd/server/server.go`: Standalone web server version.
- `parts.db`: The SQLite database containing parts and manuals.
- `cmd/builder/main.go`: The script used to parse PDFs and populate the database.
- `build_release.ps1`: PowerShell script for building and packaging releases.

## Search Logic Mandates
- Search must support partial part numbers (e.g., "WG8", "5935") and full part numbers ("WG8-5935").
- The search query should be normalized (hyphens removed) to match against a normalized version of the part numbers in the database.
- Always search across `part_number`, `base_part`, `description`, and `remarks`.

## Database Schema Notes
- The `parts` table must include a `remarks` column.
- If `remarks` is missing, it must be added via `ALTER TABLE parts ADD COLUMN remarks TEXT;`.

## Development Workflows
- Before building, run `templ generate ./bow-gui/`.
- Use `go run cmd/tools/verify_search.go` to validate search logic against the database before committing changes.
