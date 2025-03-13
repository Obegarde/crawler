package main

import (
	"fmt"
	"net/url"
)

func normalizeURL(inputURL string) (string, error) {
	URLStruct, err := url.Parse(inputURL)
	if err != nil {
		return "", fmt.Errorf("couldn't parse URL:%v", err)
	}
	return URLStruct.String(), nil
}
