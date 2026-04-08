package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "modernc.org/sqlite"
)

type PartOccurrence struct {
	ModelSeries    string
	ManualRevision string
	FigureID       string
	KeyNo          string
	FullPartNumber string
	Revision       string
	Description    string
	Remarks        string
}

type GroupedResult struct {
	BasePart    string
	Description string
	Occurrences []PartOccurrence
}

func performSmartSearch(db *sql.DB, q string) []GroupedResult {
	if q == "" {
		return nil
	}
	q = strings.ToUpper(strings.TrimSpace(q))

	// Search term for fields that might have hyphens
	likeQ := "%" + q + "%"

	// Normalized search term for fields where we strip hyphens
	normalizedQ := strings.ReplaceAll(q, "-", "")
	likeNormalizedQ := "%" + normalizedQ + "%"
	
	sqlQuery := `
		SELECT p.base_part, p.description, m.model_series, m.revision, p.figure_id, p.key_no, p.part_number, p.revision, p.remarks
		FROM parts p
		JOIN manuals m ON p.manual_id = m.id
		WHERE p.part_number LIKE ? 
		   OR REPLACE(p.part_number, '-', '') LIKE ?
		   OR p.base_part LIKE ?
		   OR REPLACE(p.base_part, '-', '') LIKE ?
		   OR p.description LIKE ?
		   OR p.remarks LIKE ?
		ORDER BY p.base_part, m.model_series
	`
	args := []interface{}{likeQ, likeNormalizedQ, likeQ, likeNormalizedQ, likeQ, likeQ}

	rows, err := db.Query(sqlQuery, args...)
	if err != nil {
		log.Printf("Search error: %v", err)
		return nil
	}
	defer rows.Close()

	var baseParts []string
	groups := make(map[string]*GroupedResult)

	for rows.Next() {
		var base, desc, model, mRev, fig, key, full, pRev, remarks string
		rows.Scan(&base, &desc, &model, &mRev, &fig, &key, &full, &pRev, &remarks)

		if _, ok := groups[base]; !ok {
			groups[base] = &GroupedResult{
				BasePart:    base,
				Description: desc,
			}
			baseParts = append(baseParts, base)
		}

		groups[base].Occurrences = append(groups[base].Occurrences, PartOccurrence{
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

	final := make([]GroupedResult, 0, len(baseParts))
	for _, b := range baseParts {
		final = append(final, *groups[b])
	}
	
	return final
}

func main() {
	db, err := sql.Open("sqlite", "file:../../parts.db?mode=ro")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	failed := false

	// 1. Verify Schema Integrity
	fmt.Printf("Verifying database schema... ")
	var schema string
	err = db.QueryRow("SELECT sql FROM sqlite_master WHERE type='table' AND name='parts'").Scan(&schema)
	if err != nil {
		fmt.Printf("FAIL (Could not read parts table schema: %v)\n", err)
		failed = true
	} else if !strings.Contains(strings.ToLower(schema), "remarks") {
		fmt.Printf("FAIL (Table 'parts' is missing 'remarks' column)\n")
		failed = true
	} else {
		fmt.Printf("PASS\n")
	}

	if failed {
		fmt.Println("\nResult: DATABASE INTEGRITY CHECK FAILED")
		os.Exit(1)
	}

	// 2. Verify Search Logic
	testCases := []struct {
		query       string
		minResults  int
		mustContain string
	}{
		{"WG8", 5, "WG8-"},
		{"5935", 1, "WG8-5935"},
		{"WG8-5935", 1, "WG8-5935"},
		{"Roller", 5, "ROLLER"},
	}

	for _, tc := range testCases {
		fmt.Printf("Testing query: '%s'... ", tc.query)
		results := performSmartSearch(db, tc.query)

		count := len(results)
		foundMatch := false
		for _, r := range results {
			if strings.Contains(r.BasePart, tc.mustContain) || strings.Contains(strings.ToUpper(r.Description), tc.mustContain) {
				foundMatch = true
				break
			}
		}

		if count < tc.minResults {
			fmt.Printf("FAIL (Expected >= %d, got %d)\n", tc.minResults, count)
			failed = true
		} else if !foundMatch {
			fmt.Printf("FAIL (Results found but none contained '%s')\n", tc.mustContain)
			failed = true
		} else {
			fmt.Printf("PASS (%d results found)\n", count)
		}
	}

	if failed {
		fmt.Println("\nResult: SOME TESTS FAILED")
		os.Exit(1)
	} else {
		fmt.Println("\nResult: ALL TESTS PASSED")
	}
}
