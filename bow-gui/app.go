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
	
	// Open database (Look for it in the same directory as the EXE)
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

// GetModels returns an HTML list of all manuals in the database
func (a *App) GetModels() string {
	if a.db == nil {
		return "Database not connected."
	}

	rows, err := a.db.Query("SELECT model_series, revision, filename FROM manuals ORDER BY model_series")
	if err != nil {
		return fmt.Sprintf("Error querying manuals: %v", err)
	}
	defer rows.Close()

	var sb strings.Builder
	sb.WriteString("<div class='bg-white p-6 rounded-xl shadow-md border border-gray-200'>")
	sb.WriteString("<h3 class='text-xl font-bold mb-4 border-b pb-2'>Manuals in Database</h3>")
	sb.WriteString("<table class='w-full text-left border-collapse text-sm'>")
	sb.WriteString("<thead class='bg-gray-50 uppercase text-xs text-gray-500 font-bold'><tr>")
	sb.WriteString("<th class='px-4 py-2 border-b'>Model Series</th>")
	sb.WriteString("<th class='px-4 py-2 border-b'>Revision</th>")
	sb.WriteString("<th class='px-4 py-2 border-b'>Source Filename</th>")
	sb.WriteString("</tr></thead><tbody>")

	count := 0
	for rows.Next() {
		var model, rev, file string
		rows.Scan(&model, &rev, &file)
		sb.WriteString("<tr class='hover:bg-gray-50'>")
		sb.WriteString(fmt.Sprintf("<td class='px-4 py-2 border-b font-bold'>%s</td>", model))
		sb.WriteString(fmt.Sprintf("<td class='px-4 py-2 border-b'><span class='bg-gray-200 px-2 py-0.5 rounded'>Rev %s</span></td>", rev))
		sb.WriteString(fmt.Sprintf("<td class='px-4 py-2 border-b text-gray-500 font-mono text-[10px]'>%s</td>", file))
		sb.WriteString("</tr>")
		count++
	}
	sb.WriteString("</tbody></table>")
	sb.WriteString(fmt.Sprintf("<p class='mt-4 text-gray-400 text-xs italic'>Total manuals: %d</p>", count))
	sb.WriteString("</div>")

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
