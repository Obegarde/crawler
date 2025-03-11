package main

import (
	"fmt"
	"strings"
)

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.mu.Lock()
	if cfg.checkedPages > cfg.maxPages {
		cfg.mu.Unlock()
		return
	} else {
		cfg.mu.Unlock()
	}
	if !strings.Contains(rawCurrentURL, cfg.baseURL.String()) {
		return
	}
	normCurrent, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error normalizing: %v/n", err)
	}
	if !cfg.addPageVisit(normCurrent) {
		return
	}
	pageHTML, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error getting page: %v\n", err)
		return
	}
	fmt.Printf("Getting HTML From :%v\n", rawCurrentURL)
	err = cfg.WriteHTMLToFile(pageHTML, normCurrent)
	if err != nil {
		fmt.Printf("Failed to save HTML: %v\n", err)
		panic(err)
	}
	links, err := getURLsFromHTML(pageHTML, rawCurrentURL)
	if err != nil {
		fmt.Printf("Error extracting URLs: %v\n", err)
		return
	}
	for _, link := range links {
		normLink, err := normalizeURL(link)
		if err != nil {
			fmt.Printf("failed to normalize: %v\n", err)
		}
		cfg.addNewPage(normLink)
		cfg.wg.Add(1)
		go func() {
			defer cfg.wg.Done()
			defer func() { <-cfg.concurrencyControl }()
			cfg.concurrencyControl <- struct{}{}
			cfg.crawlPage(link)
		}()
	}
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	if cfg.pages[normalizedURL] == 0 {
		cfg.checkedPages += 1
		cfg.pages[normalizedURL] += 1
		return true
	}
	cfg.pages[normalizedURL] += 1
	return false
}

func (cfg *config) addNewPage(normalizedURL string) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	_, ok := cfg.pages[normalizedURL]
	if ok {
		return
	}
	cfg.pages[normalizedURL] = 0
}
