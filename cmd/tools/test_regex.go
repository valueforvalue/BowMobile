package main

import (
	"fmt"
	"regexp"
)

func main() {
	partRegex := regexp.MustCompile(`^(\d+)([A-Z0-9]{3}-[A-Z0-9]{4})-([A-Z0-9]{3})(\d+)(.*)$`)
	
	testLines := []string{
		"1FE3-9169-0004CAP",
		"2FE8-6963-0001COVER, RIGHT FRONT, UP. INNER",
		"3FE8-6968-0001COVER, CONNECTOR, RIGHT, LOWER",
		"19XA9-0893-00015SCREW, M4X8",
	}

	for _, line := range testLines {
		matches := partRegex.FindStringSubmatch(line)
		if matches != nil {
			fmt.Printf("Match: %s\n", line)
			fmt.Printf("  Key: %s\n", matches[1])
			fmt.Printf("  Base: %s\n", matches[2])
			fmt.Printf("  Rev: %s\n", matches[3])
			fmt.Printf("  Qty: %s\n", matches[4])
			fmt.Printf("  Desc: %s\n", matches[5])
		} else {
			fmt.Printf("No match: %s\n", line)
		}
	}
}
