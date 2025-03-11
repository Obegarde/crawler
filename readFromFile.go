package main

import (
	"bufio"
	"encoding/gob"
	"os"
)

func ReadPagesMapFromFile(pagesPath string) (map[string]int, error) {
	file, err := os.Open(pagesPath)
	if err != nil {
		return map[string]int{}, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var pagesMap map[string]int
	decoder := gob.NewDecoder(reader)
	err = decoder.Decode(&pagesMap)
	if err != nil {
		return map[string]int{}, err
	}
	return pagesMap, nil
}
