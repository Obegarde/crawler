package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func getHTML(rawURL string) (string, error) {
	res, err := http.Get(rawURL)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		return "", fmt.Errorf("response failed with status code: %d and\nbody: %s", res.StatusCode, body)
	}
	if err != nil {
		return "", err
	}
	if len(res.Header["Content-Type"]) == 0 {
		return "", fmt.Errorf("no Content-Type header found")
	}

	if !strings.Contains(res.Header["Content-Type"][0], "text/html") {
		return "", fmt.Errorf("content type not text/html")
	}
	return string(body), nil
}
