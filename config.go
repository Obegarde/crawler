package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
)

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
	checkedPages       int
}

func configure(rawURL, maxConcurrencyControl, maxPagesString string) (*config, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return &config{}, fmt.Errorf("failed to parse url:%v", err)
	}
	maxPagesInt, err := strconv.Atoi(maxPagesString)
	if err != nil {
		return &config{}, fmt.Errorf("failed to convert maxpages:%v", err)
	}
	maxConcurrency, err := strconv.Atoi(maxConcurrencyControl)
	if err != nil {
		return &config{}, fmt.Errorf("failed to convert concurrenycontrol: %v", err)
	}
	err = os.MkdirAll("out/"+parsedURL.Host, 0750)
	if err != nil {
		return &config{}, fmt.Errorf("failed to create output directory: %v", err)
	}
	loadedPages, err := ReadPagesMapFromFile("out/" + parsedURL.Host + "pagesMap")
	if err != nil {
		fmt.Println("Failed to load pagesMap using cleanMap")
	}

	newConfig := config{
		pages:              loadedPages,
		baseURL:            parsedURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPagesInt,
		checkedPages:       0,
	}
	return &newConfig, nil
}
