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

	fmt.Println("--- Investigating potential 'junk' data ---")
	
	fmt.Println("\n1. Top 10 most common descriptions (to find repeating noise):")
	rows, _ := db.Query("SELECT description, count(*) as c FROM parts GROUP BY description ORDER BY c DESC LIMIT 10")
	for rows.Next() {
		var d string
		var c int
		rows.Scan(&d, &c)
		fmt.Printf("[%d] %q\n", c, d)
	}
	rows.Close()

	fmt.Println("\n2. Parts where description is very short (< 3 chars):")
	rows, _ = db.Query("SELECT part_number, description, figure_id FROM parts WHERE length(description) < 3 LIMIT 10")
	for rows.Next() {
		var p, d, f string
		rows.Scan(&p, &d, &f)
		fmt.Printf("Part: %s | Desc: %q | Fig: %s\n", p, d, f)
	}
	rows.Close()

	fmt.Println("\n3. Parts that look like dates (e.g. 2026-03-09):")
	rows, _ = db.Query("SELECT part_number, description FROM parts WHERE part_number LIKE '202%-%' LIMIT 10")
	for rows.Next() {
		var p, d string
		rows.Scan(&p, &d)
		fmt.Printf("Part: %s | Desc: %q\n", p, d)
	}
	rows.Close()
}
