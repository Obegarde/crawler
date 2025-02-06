package main
import(
	"strings"
	"fmt"
)

func (cfg *config)crawlPage(rawCurrentURL string){
	cfg.mu.Lock()
	if len(cfg.pages) > cfg.maxPages{
		cfg.mu.Unlock()
		return
	}else{
		cfg.mu.Unlock()	
	} 
	if !strings.Contains(rawCurrentURL,cfg.baseURL.String()){
		return 
	}
	normCurrent,err := normalizeURL(rawCurrentURL)
	if err != nil{
		fmt.Printf("Error normalizing: %v/n", err)
	}
	if !cfg.addPageVisit(normCurrent){
		return
	}
	
	pageHTML,err := getHTML(rawCurrentURL)
	if err!= nil{
		fmt.Printf("Error getting page: %v\n", err)
		return
	}
	fmt.Printf("Getting HTML From :%v\n", rawCurrentURL)

	links, err := getURLsFromHTML(pageHTML,rawCurrentURL)
	if err != nil{
		fmt.Printf("Error extracting URLs: %v\n", err)
		return
	}
	for _,link := range(links){
		cfg.wg.Add(1)
		go func(){
			defer cfg.wg.Done()
			defer func(){<-cfg.concurrencyControl}()
			cfg.concurrencyControl <- struct{}{}
			cfg.crawlPage(link)
		}()
	}
}

func (cfg *config)addPageVisit(normalizedURL string)(isFirst bool){
	cfg.mu.Lock()
	defer cfg.mu.Unlock()	
	_, ok := cfg.pages[normalizedURL]
	if ok{
		cfg.pages[normalizedURL] += 1
		return false
	}
	cfg.pages[normalizedURL] = 1
	return true
}

