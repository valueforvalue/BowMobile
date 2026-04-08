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

	"bow" // Import root package for shared types

	_ "modernc.org/sqlite"
)

func main() {
	// Reaching up to parent directory for parts.db
	db, err := sql.Open("sqlite", "parts.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = initDB(db)
	if err != nil {
		log.Fatal(err)
	}

	var filesToProcess []string
	if len(os.Args) > 1 {
		// Single file mode
		filesToProcess = append(filesToProcess, os.Args[1])
	} else {
		// Folder mode - look in parent Parts directory
		entries, err := os.ReadDir("Parts")
		if err != nil {
			log.Fatal(err)
		}
		for _, e := range entries {
			fname := e.Name()
			if strings.HasSuffix(strings.ToLower(fname), ".pdf") &&
				fname != "decrypted.pdf" &&
				fname != "temp_decrypted.pdf" {
				filesToProcess = append(filesToProcess, filepath.Join("Parts", fname))
			}
		}
	}

	totalStart := time.Now()

	for _, pdfPath := range filesToProcess {
		fname := filepath.Base(pdfPath)
		
		// Skip if already in DB
		var exists bool
		db.QueryRow("SELECT EXISTS(SELECT 1 FROM manuals WHERE filename=?)", fname).Scan(&exists)
		if exists && len(os.Args) == 1 {
			continue
		}

		manualStart := time.Now()
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

	fmt.Printf("\nProcessing complete in %v\n", time.Since(totalStart).Round(time.Second))
}

func extractManualInfo(pdfPath, filename string) bow.ManualInfo {
	info := bow.ManualInfo{ModelSeries: "Unknown", Revision: "Unknown"}
	reRev := regexp.MustCompile(`_r(\d+)_`)
	if revMatch := reRev.FindStringSubmatch(filename); len(revMatch) > 1 {
		info.Revision = revMatch[1]
	}
	reModel := regexp.MustCompile(`(?:imageFORCE|imagePRESS|imagePRESS_Lite|imageRUNNER_ADVANCE|imageRUNNER_ADVANCE_DX)_([A-Z0-9_]+)_(?:PC|Series)`)
	if m := reModel.FindStringSubmatch(filename); len(m) > 1 {
		info.ModelSeries = strings.ReplaceAll(m[1], "_", "/")
	}
	pdftotextPath := `C:\Program Files\Git\mingw64\bin\pdftotext.exe`
	cmd := exec.Command(pdftotextPath, "-l", "5", pdfPath, "-")
	output, err := cmd.Output()
	if err == nil {
		content := string(output)
		if info.Revision == "Unknown" {
			if revMatch := regexp.MustCompile(`Rev\.\s*(\d+)`).FindStringSubmatch(content); revMatch != nil {
				info.Revision = revMatch[1]
			}
		}
		if info.ModelSeries == "Unknown" {
			patterns := []string{
				`imageRUNNER\s+ADVANCE\s+DX\s+([A-Z0-9 /iF]+?)(?:\s+Series|\s+C\d+)`,
				`imageRUNNER\s+ADVANCE\s+([A-Z0-9 /iF]+?)(?:\s+Series|\s+C\d+)`,
				`imageFORCE\s+([A-Z0-9 /iF]+?)(?:\s+Series|\s+C\d+)`,
				`imagePRESS\s+([A-Z0-9 /iF]+?)(?:\s+Series|\s+C\d+)`,
			}
			for _, p := range patterns {
				if m := regexp.MustCompile(p).FindStringSubmatch(content); m != nil {
					info.ModelSeries = strings.TrimSpace(m[1])
					break
				}
			}
		}
	}
	return info
}

func processWithPDFToText(db *sql.DB, pdfPath string, manualID int64) {
	pdftotextPath := `C:\Program Files\Git\mingw64\bin\pdftotext.exe`
	cmd := exec.Command(pdftotextPath, "-table", pdfPath, "-")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error running pdftotext: %v", err)
		return
	}

	content := string(output)
	lines := strings.Split(content, "\n")
	
	// Anchor on Part Number: XXX-XXXX-XXX
	anchorRegex := regexp.MustCompile(`([A-Z0-9]{3}-[A-Z0-9]{4})-([A-Z0-9]{3})`)
	qtyRegex := regexp.MustCompile(`^(\d+)\s+(.*)$`)
	keyNoRegex := regexp.MustCompile(`^\d{1,3}$`)
	remSplitRegex := regexp.MustCompile(`^(.*?)\s{3,}(.*)$`)
	
	var currentFigureID string

	for _, line := range lines {
		line = strings.TrimRight(line, "\r\n")
		
		// 1. Noise Filter & Figure Detection
		if strings.Contains(line, "Figure") {
			fParts := strings.Split(line, "Figure")
			if len(fParts) > 1 {
				fIDMatch := regexp.MustCompile(`([A-Z0-9]+)`).FindStringSubmatch(strings.TrimSpace(fParts[1]))
				if fIDMatch != nil {
					currentFigureID = fIDMatch[1]
				}
			}
			continue
		}

		if strings.HasPrefix(currentFigureID, "SNL") || currentFigureID == "" {
			continue
		}

		if strings.Contains(line, "Page") || strings.Contains(line, "Copyright") || strings.Contains(line, "202") {
			continue
		}

		// 2. Find Part Number Anchor
		loc := anchorRegex.FindStringSubmatchIndex(line)
		if loc == nil {
			continue
		}

		// 3. Extract Fields
		basePart := line[loc[2]:loc[3]]
		revision := line[loc[4]:loc[5]]
		
		// KeyNo is to the left of the anchor
		keyNo := strings.TrimSpace(line[:loc[0]])
		
		// Description and Qty are to the right
		rightSide := strings.TrimSpace(line[loc[1]:])
		qty := ""
		description := rightSide
		remarks := ""
		
		if qr := qtyRegex.FindStringSubmatch(rightSide); len(qr) > 2 {
			qty = qr[1]
			description = qr[2]
		}

		// Split Remarks (3+ spaces)
		if remParts := remSplitRegex.FindStringSubmatch(description); len(remParts) > 2 {
			description = remParts[1]
			remarks = remParts[2]
		}

		// 4. Validation
		if !keyNoRegex.MatchString(keyNo) || len(description) < 3 {
			continue
		}

		savePart(db, manualID, currentFigureID, bow.Part{
			KeyNo:       keyNo,
			BasePart:    basePart,
			Revision:    revision,
			Qty:         qty,
			Description: description,
			Remarks:     remarks,
		})
	}
}

func initDB(db *sql.DB) error {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS manuals (id INTEGER PRIMARY KEY AUTOINCREMENT, filename TEXT UNIQUE, model_series TEXT, revision TEXT);
	CREATE TABLE IF NOT EXISTS figures (manual_id INTEGER, id TEXT, PRIMARY KEY (manual_id, id), FOREIGN KEY(manual_id) REFERENCES manuals(id));
	CREATE TABLE IF NOT EXISTS parts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		manual_id INTEGER,
		figure_id TEXT,
		key_no TEXT,
		part_number TEXT,
		base_part TEXT,
		revision TEXT,
		qty TEXT,
		description TEXT,
		remarks TEXT,
		FOREIGN KEY(manual_id) REFERENCES manuals(id)
	);
	CREATE INDEX IF NOT EXISTS idx_base_part ON parts(base_part);
	`
	_, err := db.Exec(sqlStmt)
	return err
}

func saveManual(db *sql.DB, filename, model, rev string) (int64, error) {
	_, err := db.Exec("INSERT OR IGNORE INTO manuals (filename, model_series, revision) VALUES (?, ?, ?)", filename, model, rev)
	if err != nil {
		return 0, err
	}
	
	var id int64
	err = db.QueryRow("SELECT id FROM manuals WHERE filename = ?", filename).Scan(&id)
	return id, err
}

func savePart(db *sql.DB, manualID int64, figureID string, p bow.Part) {
	db.Exec("INSERT OR IGNORE INTO figures (manual_id, id) VALUES (?, ?)", manualID, figureID)
	db.Exec("INSERT INTO parts (manual_id, figure_id, key_no, part_number, base_part, revision, qty, description, remarks) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		manualID, figureID, p.KeyNo, p.BasePart + "-" + p.Revision, p.BasePart, p.Revision, p.Qty, p.Description, p.Remarks)
}
