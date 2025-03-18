package main

import (
	"fmt"
	"net/url"
	"slices"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func getURLsFromHTML(htmlBody, rawBaseURL string) (map[string]CrawlInfo, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse base URL:%v", err)
	}

	HTMLReader := strings.NewReader(htmlBody)
	urls := make(map[string]CrawlInfo)
	doc, err := html.Parse(HTMLReader)
	if err != nil {
		return nil, err
	}
	for node := range doc.Descendants() {
		if node.Type == html.ElementNode && node.DataAtom == atom.A {
			shouldDownload := false

			for nodeDec := range node.Descendants() {
				if nodeDec.Type == html.TextNode && slices.Contains([]string{"Se jobbet"}, nodeDec.Data) {
					shouldDownload = true
					break
				}
			}
			for _, a := range node.Attr {
				if a.Key == "href" {
					hrefVal, err := url.Parse(a.Val)
					if err != nil {
						fmt.Printf("couldn't parse href '%v': %v\n", a.Val, err)
						continue
					}
					resolvedURL := baseURL.ResolveReference(hrefVal)
					urls[resolvedURL.String()] = CrawlInfo{
						Url:            resolvedURL,
						Checked:        false,
						ShouldDownload: shouldDownload,
					}
				}
			}
		}
	}
	return urls, nil
}
