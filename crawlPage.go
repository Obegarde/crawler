package main

import (
	"fmt"
	"strings"
)

func (cfg *config) crawlPage(currentURL string) {
	fmt.Println(currentURL)
	pageHTML, err := getHTML(currentURL)
	if err != nil {
		fmt.Printf("Error getting page: %v\n", err)
		return
	}
	// fmt.Printf("Getting HTML From :%v\n", currentURL)
	err = cfg.WriteHTMLToFile(pageHTML, currentURL)
	if err != nil {
		fmt.Printf("Failed to save HTML: %v\n", err)
		panic(err)
	}

	newLinks, err := getURLsFromHTML(pageHTML, currentURL)
	if err != nil {
		fmt.Printf("Error extracting URLs: %v\n", err)
		return
	}
	for _, newLink := range newLinks {
		normLink, err := normalizeURL(newLink)
		if err != nil {
			fmt.Printf("Failed to normalizeURL: %v", err)
			continue
		}
		cfg.addNewPage(normLink)
	}
}

func (cfg *config) workThroughLinks(links []string) int {
	linksWorkedThrough := 0
	for _, link := range links {
		linksWorkedThrough += 1
		if linksWorkedThrough > cfg.maxPages {
			break
		}
		cfg.setPageChecked(link)
		cfg.wg.Add(1)
		go func() {
			defer cfg.wg.Done()
			defer func() { <-cfg.concurrencyControl }()
			cfg.concurrencyControl <- struct{}{}
			cfg.crawlPage(link)
		}()
	}
	return linksWorkedThrough
}

func (cfg *config) generateLinkList() ([]string, error) {
	links := []string{}
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	// fmt.Println(cfg.pages)
	for key, val := range cfg.pages {
		if !val {
			links = append(links, key)
		}
	}
	if len(links) == 0 {
		return links, fmt.Errorf("no unchecked links found")
	}
	return links, nil
}

func (cfg *config) setPageChecked(normalizedURL string) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	cfg.pages[normalizedURL] = true
}

func (cfg *config) addNewPage(normalizedURL string) {
	if !strings.Contains(normalizedURL, cfg.baseURL.Host) {
		return
	}
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	_, ok := cfg.pages[normalizedURL]
	if ok {
		return
	}
	cfg.pages[normalizedURL] = false
}
