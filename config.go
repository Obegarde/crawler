package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
)

type CrawlInfo struct {
	Url            *url.URL
	Checked        bool
	ShouldDownload bool
}

type config struct {
	pages              map[string]CrawlInfo
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
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
	loadedPages, err := ReadPagesMapFromFile("out/" + parsedURL.Host + "/pagesMap")
	if err != nil {
		fmt.Printf("Failed to load pagesMap using cleanMap: %v\n", err)
	}

	newConfig := config{
		pages:              loadedPages,
		baseURL:            parsedURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPagesInt,
	}
	return &newConfig, nil
}
