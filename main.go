package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	argsWithProg := os.Args
	if argsWithProg[1] == "help" {
		fmt.Println("This is a webcrawler that saves the pages it visits at out/")
		fmt.Println("Example initialization:")
		fmt.Println(`./crawler "https://crawler-test.com/" 3 25`)
		fmt.Println("Crawler itself, Raw url, amount of concurrent threads, pages to visit before stopping")
		return
	}
	if len(argsWithProg) < 4 {
		fmt.Println("Must provide both website, concurrent threads and max pages to crawl")
		os.Exit(1)
	} else if len(argsWithProg) > 4 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	} else {
		cfg, err := configure(argsWithProg[1], argsWithProg[2], argsWithProg[3])
		if err != nil {
			fmt.Printf("Failed to configure crawler:%v\n", err)
			os.Exit(1)
		}
		defer cfg.WritePagesMapToFile("pagesMap")

		normalizedBase, err := normalizeURL(cfg.baseURL.String())
		if err != nil {
			fmt.Printf("failed to format baseURL: %v", err)
			return
		}
		baseURLStruct := CrawlInfo{
			URL:            normalizedBase,
			Checked:        false,
			ShouldDownload: false,
		}
		cfg.addNewPage(normalizedBase, baseURLStruct)
		fmt.Println("starting crawl")
		ctx, cancel := context.WithCancel(context.Background())
		go func(ctx context.Context) {
			checkedPages := 0
			for checkedPages < cfg.maxPages {
				links, err := cfg.generateLinkList()
				if err != nil {
					fmt.Printf("failed to generateLinkList: %v\n", err)
					break
				}
				checkedPages += cfg.workThroughLinks(links, ctx)
				cfg.wg.Wait()
				fmt.Println("Crawl finished: Press Ctrl + c to stop program")
			}
		}(ctx)
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
		<-signalCh
		fmt.Println("stopping Crawler")
		cancel()
		printReport(cfg.pages, cfg.baseURL.String())
	}
}
