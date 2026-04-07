package main

import (
	"fmt"
	"log"
	"time"

	"github.com/dslipak/pdf"
)

func main() {
	pdfPath := "Parts/temp_decrypted.pdf"
	r, err := pdf.Open(pdfPath)
	if err != nil {
		log.Fatal(err)
	}

	numPages := r.NumPage()

	for i := 50; i <= 100 && i <= numPages; i++ {
		fmt.Printf("Testing Page %d... ", i)
		p := r.Page(i)
		
		done := make(chan bool)
		go func() {
			p.GetPlainText(nil)
			done <- true
		}()

		select {
		case <-done:
			fmt.Println("OK")
		case <-time.After(10 * time.Second):
			fmt.Println("FAIL (Hanging)")
			return
		}
	}
}
