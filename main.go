package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type Part struct {
	KeyNo       string
	BasePart    string
	Revision    string
	Qty         string
	Description string
}

type ManualInfo struct {
	ModelSeries string
	Revision    string
}

func main() {
	db, err := sql.Open("sqlite", "parts.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = initDB(db)
	if err != nil {
		log.Fatal(err)
	}

	files, err := os.ReadDir("Parts")
	if err != nil {
		log.Fatal(err)
	}

	totalStart := time.Now()

	for _, file := range files {
		fname := file.Name()
		if !strings.HasSuffix(strings.ToLower(fname), ".pdf") || 
		   fname == "decrypted.pdf" || 
		   fname == "temp_decrypted.pdf" {
			continue
		}

		manualStart := time.Now()
		pdfPath := filepath.Join("Parts", fname)
		fmt.Printf("\n--- Starting: %s ---\n", fname)

		tempPath := filepath.Join("Parts", "temp_decrypted.pdf")
		os.Remove(tempPath)

		qpdfPath := `C:\Program Files\qpdf 12.3.2\bin\qpdf.exe`
		fmt.Print("Step 1: Decrypting... ")
		cmd := exec.Command(qpdfPath, "--decrypt", pdfPath, tempPath)
		if err := cmd.Run(); err != nil {
			log.Printf("Error decrypting %s: %v", fname, err)
			continue
		}
		fmt.Printf("Done (%v)\n", time.Since(manualStart).Round(time.Millisecond))

		// Extract Model and Revision
		info := extractManualInfo(tempPath, fname)
		fmt.Printf("Detected: %s (Rev %s)\n", info.ModelSeries, info.Revision)

		manualID, err := saveManual(db, fname, info.ModelSeries, info.Revision)
		if err != nil {
			log.Printf("Error saving manual %s: %v", fname, err)
			os.Remove(tempPath)
			continue
		}

		fmt.Println("Step 2: Parsing with pdftotext...")
		processWithPDFToText(db, tempPath, manualID)
		os.Remove(tempPath)
		fmt.Printf("--- Finished: %s (Time: %v) ---\n", fname, time.Since(manualStart).Round(time.Second))
	}

	fmt.Printf("\nAll manuals processed in %v\n", time.Since(totalStart).Round(time.Second))
}

func extractManualInfo(pdfPath, filename string) ManualInfo {
	info := ManualInfo{ModelSeries: "Unknown", Revision: "Unknown"}

	// 1. More aggressive filename extraction
	// We look for anything between DX_ and the next marker (_Series or _PC)
	reFilename := regexp.MustCompile(`DX_([^_]+(?:_[^_]+)*?)_(?:Series|PC)`)
	fMatches := reFilename.FindStringSubmatch(filename)
	if len(fMatches) > 1 {
		info.ModelSeries = strings.ReplaceAll(fMatches[1], "_", "/")
	}
	
	reRev := regexp.MustCompile(`_r(\d+)_`)
	revMatch := reRev.FindStringSubmatch(filename)
	if len(revMatch) > 1 {
		info.Revision = revMatch[1]
	}

	// 2. First Page content extraction (Secondary/Backup)
	pdftotextPath := `C:\Program Files\Git\mingw64\bin\pdftotext.exe`
	// Read first 5 pages to be safe, sometimes series info is on page 2 or 3
	cmd := exec.Command(pdftotextPath, "-l", "5", pdfPath, "-")
	output, err := cmd.Output()
	if err == nil {
		content := string(output)
		
		if info.Revision == "Unknown" {
			revMatch := regexp.MustCompile(`Rev\.\s*(\d+)`).FindStringSubmatch(content)
			if revMatch != nil {
				info.Revision = revMatch[1]
			}
		}

		if info.ModelSeries == "Unknown" || info.ModelSeries == "" {
			// Look for "imageRUNNER ADVANCE DX [MODEL]"
			seriesMatch := regexp.MustCompile(`DX\s+([A-Z0-9 /iF]+?)(?:\s+Series|\s+C\d+)`).FindStringSubmatch(content)
			if seriesMatch != nil {
				info.ModelSeries = strings.TrimSpace(seriesMatch[1])
			}
		}
	}

	return info
}

func processWithPDFToText(db *sql.DB, pdfPath string, manualID int64) {
	pdftotextPath := `C:\Program Files\Git\mingw64\bin\pdftotext.exe`
	cmd := exec.Command(pdftotextPath, "-layout", pdfPath, "-")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error running pdftotext: %v", err)
		return
	}

	content := string(output)
	lines := strings.Split(content, "\n")
	
	partRegex := regexp.MustCompile(`(\d+)\s+([A-Z0-9]{3}-[A-Z0-9]{4})-([A-Z0-9]{3})\s+(\d*)\s*(.*)`)
	var currentFigureID string

	for _, line := range lines {
		line = strings.TrimRight(line, "\r\n")
		
		if strings.Contains(line, "Figure") {
			fParts := strings.Split(line, "Figure")
			if len(fParts) > 1 {
				fIDMatch := regexp.MustCompile(`([A-Z0-9]+)`).FindStringSubmatch(strings.TrimSpace(fParts[1]))
				if fIDMatch != nil {
					currentFigureID = fIDMatch[1]
				}
			}
		}

		if currentFigureID != "" {
			matches := partRegex.FindAllStringSubmatch(line, -1)
			for _, m := range matches {
				part := Part{
					KeyNo:       m[1],
					BasePart:    m[2],
					Revision:    m[3],
					Qty:         m[4],
					Description: strings.TrimSpace(m[5]),
				}
				savePart(db, manualID, currentFigureID, part)
			}
		}
	}
}

func initDB(db *sql.DB) error {
	db.Exec("DROP TABLE IF EXISTS parts")
	db.Exec("DROP TABLE IF EXISTS figures")
	db.Exec("DROP TABLE IF EXISTS manuals")
	
	sqlStmt := `
	CREATE TABLE manuals (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		filename TEXT,
		model_series TEXT,
		revision TEXT
	);
	CREATE TABLE figures (
		manual_id INTEGER,
		id TEXT,
		PRIMARY KEY (manual_id, id),
		FOREIGN KEY(manual_id) REFERENCES manuals(id)
	);
	CREATE TABLE parts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		manual_id INTEGER,
		figure_id TEXT,
		key_no TEXT,
		part_number TEXT,
		base_part TEXT,
		revision TEXT,
		qty TEXT,
		description TEXT,
		FOREIGN KEY(manual_id) REFERENCES manuals(id)
	);
	CREATE INDEX idx_base_part ON parts(base_part);
	`
	_, err := db.Exec(sqlStmt)
	return err
}

func saveManual(db *sql.DB, filename, model, rev string) (int64, error) {
	res, err := db.Exec("INSERT INTO manuals (filename, model_series, revision) VALUES (?, ?, ?)", filename, model, rev)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func savePart(db *sql.DB, manualID int64, figureID string, p Part) {
	db.Exec("INSERT OR IGNORE INTO figures (manual_id, id) VALUES (?, ?)", manualID, figureID)
	db.Exec("INSERT INTO parts (manual_id, figure_id, key_no, part_number, base_part, revision, qty, description) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		manualID, figureID, p.KeyNo, p.BasePart + "-" + p.Revision, p.BasePart, p.Revision, p.Qty, p.Description)
}
