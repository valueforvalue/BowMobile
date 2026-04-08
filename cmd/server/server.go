package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"bow" // Import root package for types

	_ "modernc.org/sqlite"
)

var db *sql.DB

func main() {
	var err error
	// Open with read-only and shared cache mode for better concurrency
	// Note: parts.db is in the parent root directory
	db, err = sql.Open("sqlite", "file:../../parts.db?mode=ro&_journal_mode=WAL")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/search", handleSearch)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("../../assets"))))

	fmt.Println("Server starting at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
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
