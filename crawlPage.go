package main

import (
	"context"
	"fmt"
	"strings"
)

func (cfg *config) crawlPage(currentInfo CrawlInfo) {
	fmt.Println(currentInfo.URL)
	pageHTML, err := getHTML(currentInfo.URL)
	if err != nil {
		fmt.Printf("Error getting page: %v\n", err)
		return
	}
	if currentInfo.ShouldDownload && !currentInfo.Checked {
		fmt.Println("Downloading page")
		cfg.WriteHTMLToFile(pageHTML, currentInfo.URL)
	}
	if !strings.Contains(currentInfo.URL, cfg.baseURL.Host) {
		return
	}
	err = cfg.getURLsFromHTML(pageHTML, currentInfo.URL)
	if err != nil {
		fmt.Printf("Error extracting URLs: %v\n", err)
		return
	}
}

func (cfg *config) workThroughLinks(links []CrawlInfo, ctx context.Context) int {
	linksWorkedThrough := 0
	for _, link := range links {
		linksWorkedThrough += 1
		if linksWorkedThrough > cfg.maxPages {
			break
		}
		cfg.wg.Add(1)
		go func(ctx context.Context) {
			select {
			case <-ctx.Done():
				fmt.Println("Cancelling Crawl")
				return
			default:
				defer cfg.wg.Done()
				defer func() { <-cfg.concurrencyControl }()
				cfg.concurrencyControl <- struct{}{}
				cfg.crawlPage(link)
				cfg.setPageChecked(link.URL)
			}
		}(ctx)
	}
	return linksWorkedThrough
}

func (cfg *config) generateLinkList() ([]CrawlInfo, error) {
	links := []CrawlInfo{}
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	// fmt.Println(cfg.pages)
	for _, val := range cfg.pages {
		if !val.Checked {
			links = append(links, val)
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
