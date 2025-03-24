package main

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func (cfg *config) getURLsFromHTML(htmlBody, rawBaseURL string) error {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return fmt.Errorf("couldn't parse base URL:%v", err)
	}

	HTMLReader := strings.NewReader(htmlBody)
	doc, err := html.Parse(HTMLReader)
	if err != nil {
		return err
	}
	for node := range doc.Descendants() {
		if node.Type == html.ElementNode && node.DataAtom == atom.A {
			for _, a := range node.Attr {
				if a.Key == "href" {
					hrefVal, err := url.Parse(a.Val)
					if err != nil {
						fmt.Printf("couldn't parse href '%v': %v\n", a.Val, err)
						continue
					}
					resolvedURLString := baseURL.ResolveReference(hrefVal).String()
					if node.Parent.Data == "h4" {
						newInfo := CrawlInfo{
							URL:            resolvedURLString,
							Checked:        false,
							ShouldDownload: true,
						}
						cfg.addNewPage(resolvedURLString, newInfo)
						cfg.crawlPage(newInfo)
						cfg.setPageChecked(resolvedURLString)

					} else {
						newInfo := CrawlInfo{
							URL:            resolvedURLString,
							Checked:        false,
							ShouldDownload: false,
						}
						cfg.addNewPage(resolvedURLString, newInfo)
					}
				}
			}
		}
	}
	return nil
}
