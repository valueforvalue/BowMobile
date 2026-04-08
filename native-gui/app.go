package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strings"

	_ "modernc.org/sqlite"
)

// App struct
type App struct {
	ctx context.Context
	db  *sql.DB
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	
	// Open database
	var err error
	a.db, err = sql.Open("sqlite", "file:parts.db?mode=ro&_journal_mode=WAL")
	if err != nil {
		log.Printf("Failed to open database: %v", err)
	}
}

func (a *App) shutdown(ctx context.Context) {
	if a.db != nil {
		a.db.Close()
	}
}

// Search returns the rendered HTML results for a query
func (a *App) Search(q string) string {
	results := a.performSmartSearch(q)
	
	// We'll use a string builder to capture templ output
	var sb strings.Builder
	err := Results(results, q).Render(a.ctx, &sb)
	if err != nil {
		return fmt.Sprintf("<div class='text-red-500'>Error rendering results: %v</div>", err)
	}
	
	return sb.String()
}

func (a *App) performSmartSearch(q string) []GroupedResult {
	if q == "" || a.db == nil {
		return nil
	}
	q = strings.ToUpper(strings.TrimSpace(q))
	
	normalized := strings.ReplaceAll(q, "-", "")
	normalized = strings.ReplaceAll(normalized, " ", "")

	var sqlQuery string
	var args []interface{}

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

	rows, err := a.db.Query(sqlQuery, args...)
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
