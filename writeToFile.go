package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"os"
)

func (cfg *config) WriteHTMLToFile(htmlString, normalizedURL string) error {
	h := sha256.New()
	h.Write([]byte(normalizedURL))
	hashedURL := h.Sum(nil)
	newFilePath := fmt.Sprintf("out/%v/%x", cfg.baseURL.Host, hashedURL)
	file, err := os.Create(newFilePath)
	if err != nil {
		return fmt.Errorf("error on creation %v", err)
	}
	defer file.Close()
	_, err = file.WriteString(normalizedURL + "\n" + htmlString)
	if err != nil {
		return fmt.Errorf("error on writingString: %v", err)
	}
	file.Sync()

	return nil
}

func (cfg *config) WritePagesMapToFile(pagesName string) error {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(cfg.pages)
	if err != nil {
		return fmt.Errorf("failed to encode map to gob: %v", err)
	}
	newFilePath := fmt.Sprintf("out/%v/%v", cfg.baseURL.Host, pagesName)
	file, err := os.Create(newFilePath)
	if err != nil {
		return fmt.Errorf("failed to create or truncate pagesMap: %v", err)
	}
	defer file.Close()
	wBytes, err := file.Write(buffer.Bytes())
	if err != nil {
		return fmt.Errorf("failed to write buffer to file: %v", err)
	}
	fmt.Printf("Succesfully saved pagesMap, %v bytes written to pagesMap\n", wBytes)
	return nil
}
