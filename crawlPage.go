package main

import (
	"fmt"
)

func (cfg *config) crawlPage(currentURL string) {
	fmt.Println(currentURL)
	pageHTML, err := getHTML(currentURL)
	if err != nil {
		fmt.Printf("Error getting page: %v\n", err)
		return
	}
	err = cfg.CheckAndSaveHTML(pageHTML, currentURL)
	if err != nil {
		fmt.Printf("failed to check and save html:%v\n", err)
	}
	newLinks, err := getURLsFromHTML(pageHTML, currentURL)
	if err != nil {
		fmt.Printf("Error extracting URLs: %v\n", err)
		return
	}
	for urlString, linkInfo := range newLinks {
		cfg.addNewPage(urlString, linkInfo)
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
		if !val.Checked {
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
	if entry, ok := cfg.pages[normalizedURL]; ok {
		entry.Checked = true
		cfg.pages[normalizedURL] = entry
	}
}

func (cfg *config) addNewPage(normalizedURL string, crawlInfo CrawlInfo) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	// If url exists already just return
	_, ok := cfg.pages[normalizedURL]
	if ok {
		return
	}
	cfg.pages[normalizedURL] = crawlInfo
}
