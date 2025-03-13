package main

import (
	"fmt"
)

type page struct {
	url   string
	count int
}

func printReport(pages map[string]bool, baseURL string) {
	uniqueLinks := 0
	checkedLinks := 0
	for _, val := range pages {
		uniqueLinks += 1
		if val {
			checkedLinks += 1
		}
	}
	fmt.Println("=============================")
	fmt.Printf("REPORT for %v\n", baseURL)
	fmt.Println("=============================")
	fmt.Printf("Found %v unique links of those %v have been checked\n", uniqueLinks, checkedLinks)
}
