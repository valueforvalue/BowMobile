package main

import (
	"fmt"
	"log"

	"github.com/dslipak/pdf"
)

func main() {
	r, err := pdf.Open("Parts/decrypted.pdf")
	if err != nil {
		log.Fatal(err)
	}
	// Note: pdf.Reader from dslipak/pdf doesn't have a Close method for this case,
	// because it doesn't open a file handle itself in the same way.
	// Actually, let's check its API. it seems to be different.

	numPages := r.NumPage()
	fmt.Printf("Total pages: %d\n", numPages)

	for i := 20; i <= 25 && i <= numPages; i++ {
		p := r.Page(i)
		content, err := p.GetPlainText(nil)
		if err != nil {
			log.Printf("Error reading page %d: %v", i, err)
			continue
		}
		fmt.Printf("--- Page %d ---\n", i)
		fmt.Println(content)
	}
}
