package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"

	_ "modernc.org/sqlite"
)

type PartOccurrence struct {
	ModelSeries     string
	ManualRevision  string
	FigureID        string
	KeyNo           string
	FullPartNumber  string
	Revision        string
	Description     string
	Remarks         string
}

type GroupedResult struct {
	BasePart    string
	Description string
	Occurrences []PartOccurrence
}

var db *sql.DB

func main() {
	var err error
	// Open with read-only and shared cache mode for better concurrency
	db, err = sql.Open("sqlite", "file:parts.db?mode=ro&_journal_mode=WAL")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/search", handleSearch)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	fmt.Println("Server starting at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", nil)
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	results := performSmartSearch(query)
	renderTemplate(w, "index", map[string]interface{}{
		"Query":   query,
		"Results": results,
	})
}

func performSmartSearch(q string) []GroupedResult {
	q = strings.ToUpper(strings.TrimSpace(q))
	
	// Normalize part number: remove spaces and hyphens for matching
	normalized := strings.ReplaceAll(q, "-", "")
	normalized = strings.ReplaceAll(normalized, " ", "")

	var sqlQuery string
	var args []interface{}

	// If it looks like a part number (starts with letter/digit, has certain length)
	isPartNumber := regexp.MustCompile(`^[A-Z0-9]{3,15}$`).MatchString(normalized)

	if isPartNumber {
		// Match against base_part or full part_number (normalized) or remarks
		// We'll use a LIKE match but try to be smart
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
		// Keyword search in description or remarks
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

	groups := make(map[string]*GroupedResult)
	var baseParts []string // To maintain order

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

	var finalResults []GroupedResult
	for _, b := range baseParts {
		finalResults = append(finalResults, *groups[b])
	}

	return finalResults
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	tmpl, err := template.New("base").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Bow - Parts Cross-Reference</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; line-height: 1.6; color: #333; max-width: 1200px; margin: 0 auto; padding: 20px; background-color: #f4f7f6; }
        .header { text-align: center; margin-bottom: 40px; }
        .logo { max-width: 400px; height: auto; margin-bottom: 20px; }
        .search-box { background: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); margin-bottom: 30px; }
        input[type="text"] { width: 70%; padding: 12px; border: 1px solid #ddd; border-radius: 4px; font-size: 16px; }
        button { padding: 12px 24px; background-color: #007bff; color: white; border: none; border-radius: 4px; cursor: pointer; font-size: 16px; }
        button:hover { background-color: #0056b3; }
        .result-card { background: white; margin-bottom: 25px; border-radius: 8px; box-shadow: 0 2px 5px rgba(0,0,0,0.05); overflow: hidden; }
        .result-header { background: #e9ecef; padding: 15px 20px; border-bottom: 1px solid #dee2e6; }
        .result-header h2 { margin: 0; font-size: 1.25rem; color: #495057; }
        .result-header span { color: #6c757d; font-weight: normal; margin-left: 10px; }
        table { width: 100%; border-collapse: collapse; }
        th, td { text-align: left; padding: 12px 20px; border-bottom: 1px solid #eee; }
        th { background-color: #f8f9fa; color: #6c757d; font-weight: 600; text-transform: uppercase; font-size: 0.75rem; letter-spacing: 1px; }
        tr:last-child td { border-bottom: none; }
        .no-results { text-align: center; padding: 40px; color: #6c757d; }
        .badge { display: inline-block; padding: 0.25em 0.4em; font-size: 75%; font-weight: 700; line-height: 1; text-align: center; white-space: nowrap; vertical-align: baseline; border-radius: 0.25rem; background-color: #6c757d; color: white; }
    </style>
</head>
<body>
    <div class="header">
        <a href="/"><img src="/assets/logo.png" alt="Bow Logo" class="logo"></a>
    </div>

    <div class="search-box">
        <form action="/search" method="get">
            <input type="text" name="q" value="{{.Query}}" placeholder="Search by part number (e.g. FE3-1525) or description (e.g. Roller)..." autofocus>
            <button type="submit">Search</button>
        </form>
    </div>

    <div id="results">
        {{if .Results}}
            {{range .Results}}
                <div class="result-card">
                    <div class="result-header">
                        <h2>{{.BasePart}} <span>{{.Description}}</span></h2>
                    </div>
                    <table>
                        <thead>
                            <tr>
                                <th>Model / Series</th>
                                <th>Catalog Rev</th>
                                <th>Figure</th>
                                <th>Key No</th>
                                <th>Full Part Number</th>
                                <th>Remarks</th>
                            </tr>
                        </thead>
                        <tbody>
                            {{range .Occurrences}}
                                <tr>
                                    <td><strong>{{.ModelSeries}}</strong></td>
                                    <td><span class="badge">Rev {{.ManualRevision}}</span></td>
                                    <td>{{.FigureID}}</td>
                                    <td>{{.KeyNo}}</td>
                                    <td><code>{{.FullPartNumber}}</code></td>
                                    <td>{{.Remarks}}</td>
                                </tr>
                            {{end}}
                        </tbody>
                    </table>
                </div>
            {{end}}
        {{else if .Query}}
            <div class="no-results">
                No parts found matching "{{.Query}}".
            </div>
        {{end}}
    </div>
</body>
</html>
`