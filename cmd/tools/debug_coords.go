package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/dslipak/pdf"
)

func main() {
	r, err := pdf.Open("Parts/decrypted.pdf")
	if err != nil {
		log.Fatal(err)
	}

	p := r.Page(25)
	objs := p.Content().Text
	
	// Sort by Y descending, then X ascending
	sort.Slice(objs, func(i, j int) bool {
		if objs[i].Y != objs[j].Y {
			return objs[i].Y > objs[j].Y
		}
		return objs[i].X < objs[j].X
	})

	for _, o := range objs {
		fmt.Printf("X: %6.2f Y: %6.2f S: %q\n", o.X, o.Y, o.S)
	}
}
