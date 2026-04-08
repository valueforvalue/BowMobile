package main

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"

	"bow" // Import root package for types

	_ "modernc.org/sqlite"
)

//go:embed assets/*
var embeddedAssets embed.FS

//go:embed parts.db
var embeddedDB []byte

var db *sql.DB
var lastPulse = time.Now()

func main() {
	// 1. Ensure database exists on disk (SQLite needs a physical file)
	dbPath := "parts.db"
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		fmt.Println("Extracting embedded database...")
		err = os.WriteFile(dbPath, embeddedDB, 0644)
		if err != nil {
			log.Fatalf("Failed to extract database: %v", err)
		}
	}

	// 2. Open database
	var err error
	db, err = sql.Open("sqlite", "file:parts.db?mode=ro&_journal_mode=WAL")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 3. Setup Routes
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/search", handleSearch)
	http.HandleFunc("/pulse", handlePulse)

	// Serve embedded assets
	assetSub, _ := fs.Sub(embeddedAssets, "assets")
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(assetSub))))

	// 4. Start Server & Open Browser
	url := "http://localhost:8080"
	fmt.Printf("Bow Server starting at %s\n", url)
	
	// Shutdown watcher: If no pulse for 30 seconds, exit.
	go func() {
		for {
			time.Sleep(5 * time.Second)
			if time.Since(lastPulse) > 30*time.Second {
				fmt.Println("No browser activity detected for 30s. Shutting down...")
				os.Exit(0)
			}
		}
	}()

	go func() {
		// Wait a second for server to initialize
		time.Sleep(1 * time.Second)
		openBrowser(url)
	}()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handlePulse(w http.ResponseWriter, r *http.Request) {
	lastPulse = time.Now()
	w.WriteHeader(http.StatusOK)
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		// Try to open Chrome explicitly first
		err = exec.Command("cmd", "/c", "start", "chrome", url).Start()
		if err != nil {
			// Fallback to default browser
			err = exec.Command("cmd", "/c", "start", url).Start()
		}
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Printf("Could not open browser: %v", err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	Layout("", nil).Render(r.Context(), w)
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	results := performSmartSearch(query)

	// If HTMX request, return only the results partial
	if r.Header.Get("HX-Request") == "true" {
		Results(results, query).Render(r.Context(), w)
		return
	}

	// Otherwise return the full layout
	Layout(query, results).Render(r.Context(), w)
}

func performSmartSearch(q string) []bow.GroupedResult {
	if q == "" {
		return nil
	}
	q = strings.ToUpper(strings.TrimSpace(q))
	
	// Normalize part number: remove spaces and hyphens for matching
	normalized := strings.ReplaceAll(q, "-", "")
	normalized = strings.ReplaceAll(normalized, " ", "")

	var sqlQuery string
	var args []interface{}

	// If it looks like a part number (starts with letter/digit, has certain length)
	isPartNumber := regexp.MustCompile(`^[A-Z0-9]{3,15}$`).MatchString(normalized)

	if isPartNumber {
		sqlQuery = `
			SELECT p.base_part, p.description, m.model_series, m.revision, p.figure_id, p.key_no, p.part_number, p.revision, p.remarks
			FROM parts p
			JOIN manuals m ON p.manual_id = m.id
			WHERE REPLACE(p.part_number, '-', '') LIKE ? 
			   OR REPLACE(p.base_part, '-', '') LIKE ?
			   OR p.remarks LIKE ?
			ORDER BY p.base_part, m.model_series
		`
		args = append(args, normalized+"%", normalized+"%", "%"+q+"%")
	} else {
		sqlQuery = `
			SELECT p.base_part, p.description, m.model_series, m.revision, p.figure_id, p.key_no, p.part_number, p.revision, p.remarks
			FROM parts p
			JOIN manuals m ON p.manual_id = m.id
			WHERE p.description LIKE ?
			   OR p.remarks LIKE ?
			ORDER BY p.base_part, m.model_series
		`
		args = append(args, "%"+q+"%", "%"+q+"%")
	}

	rows, err := db.Query(sqlQuery, args...)
	if err != nil {
		log.Printf("Search error: %v", err)
		return nil
	}
	defer rows.Close()

	groups := make(map[string]*bow.GroupedResult)
	var baseParts []string // To maintain order

	for rows.Next() {
		var base, desc, model, mRev, fig, key, full, pRev, remarks string
		rows.Scan(&base, &desc, &model, &mRev, &fig, &key, &full, &pRev, &remarks)

		if _, ok := groups[base]; !ok {
			groups[base] = &bow.GroupedResult{
				BasePart:    base,
				Description: desc,
			}
			baseParts = append(baseParts, base)
		}

		groups[base].Occurrences = append(groups[base].Occurrences, bow.PartOccurrence{
			ModelSeries:    model,
			ManualRevision: mRev,
			FigureID:       fig,
			KeyNo:          key,
			FullPartNumber: full,
			Revision:       pRev,
			Description:    desc,
			Remarks:        remarks,
		})
	}

	var finalResults []bow.GroupedResult
	for _, b := range baseParts {
		finalResults = append(finalResults, *groups[b])
	}

	return finalResults
}
