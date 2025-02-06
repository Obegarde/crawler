package main
import(
	"net/url"
	"sync"
	"strconv"
	"fmt"
)
type config struct{
	pages map[string]int
	baseURL	*url.URL
	mu *sync.Mutex
	concurrencyControl chan struct{}
	wg	*sync.WaitGroup
	maxPages int
}

func configure(rawURL,maxConcurrencyControl,maxPagesString string)(*config,error){	
		parsedURL, err := url.Parse(rawURL)
		if err != nil{
		return &config{}, fmt.Errorf("failed to parse url:%v",err)
		}
		maxPagesInt,err := strconv.Atoi(maxPagesString)
		if err != nil{
		return &config{}, fmt.Errorf("failed to convert maxpages:%v",err)
	}
		maxConcurrency, err := strconv.Atoi(maxConcurrencyControl)
		if err != nil{
			return &config{}, fmt.Errorf("failed to convert concurrenycontrol: %v", err)
		}
	

		newConfig := config{
			pages: make(map[string]int),
			baseURL: parsedURL,
			mu: &sync.Mutex{},
			concurrencyControl: make(chan struct{},maxConcurrency),
			wg: &sync.WaitGroup{},
			maxPages: maxPagesInt,
		}
	return &newConfig, nil
}
