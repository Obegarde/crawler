package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"os"
)

func ReadPagesMapFromFile(pagesPath string) (map[string]bool, error) {
	pagesMap := make(map[string]bool)
	file, err := os.Open(pagesPath)
	if err != nil {
		return pagesMap, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	decoder := gob.NewDecoder(reader)
	err = decoder.Decode(&pagesMap)
	if err != nil {
		return pagesMap, err
	}
	fmt.Println("Read from pagesMap")
	return pagesMap, nil
}
