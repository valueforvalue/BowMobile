package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "parts.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var count int
	db.QueryRow("SELECT count(*) FROM parts").Scan(&count)
	fmt.Printf("Total parts in DB: %d\n", count)

	fmt.Println("\n--- Manuals in Database ---")
	rows, err := db.Query("SELECT id, model_series, revision, filename FROM manuals")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Printf("%-5s | %-15s | %-10s | %-30s\n", "ID", "Series", "Revision", "Filename")
	fmt.Println(string(make([]byte, 70)))
	for rows.Next() {
		var id int
		var series, rev, fname string
		rows.Scan(&id, &series, &rev, &fname)
		fmt.Printf("%-5d | %-15s | %-10s | %-30s\n", id, series, rev, fname)
	}

	fmt.Println("\n--- Sample Cross-Reference Matches ---")
	query := `
		SELECT p1.base_part, p1.description, m1.model_series, m2.model_series
		FROM parts p1
		JOIN parts p2 ON p1.base_part = p2.base_part AND p1.manual_id < p2.manual_id
		JOIN manuals m1 ON p1.manual_id = m1.id
		JOIN manuals m2 ON p2.manual_id = m2.id
		GROUP BY p1.base_part
		LIMIT 10
	`
	rows2, _ := db.Query(query)
	defer rows2.Close()
	fmt.Printf("%-12s | %-30s | %-12s | %-12s\n", "Base Part", "Description", "Model A", "Model B")
	for rows2.Next() {
		var base, desc, m1, m2 string
		rows2.Scan(&base, &desc, &m1, &m2)
		fmt.Printf("%-12s | %-30s | %-12s | %-12s\n", base, desc, m1, m2)
	}
}
