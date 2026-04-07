package main

import (
	"fmt"
	"log"
	"time"

	"github.com/dslipak/pdf"
)

func main() {
	// First decrypt if not already there
	pdfPath := "Parts/temp_decrypted.pdf"
	
	fmt.Println("Attempting to parse Page 50 of C3800...")
	r, err := pdf.Open(pdfPath)
	if err != nil {
		log.Fatal(err)
	}

	p := r.Page(50)
	
	done := make(chan bool)
	go func() {
		content, _ := p.GetPlainText(nil)
		fmt.Printf("Content length: %d\n", len(content))
		fmt.Println("First 100 chars:", content[:100])
		done <- true
	}()

	select {
	case <-done:
		fmt.Println("Page 40 parsed successfully.")
	case <-time.After(30 * time.Second):
		fmt.Println("TIMEOUT: Page 40 took more than 30 seconds. Potential infinite loop in library.")
	}
}
