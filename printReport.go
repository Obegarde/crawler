package main
import(
	"fmt"
	"sort"
)

type page struct{
	url string
	count int
}

func printReport(pages map[string]int,baseURL string){
	pageSlice := make([]page,len(pages))
	i := 0
	for url,value := range pages {
		pageSlice[i].url = url
		pageSlice[i].count = value
		i++
	}

	sort.SliceStable(pageSlice, func(i, j int) bool {
		if pageSlice[i].count != pageSlice[j].count{
			return pageSlice[i].count > pageSlice[j].count
		}else{
			return pageSlice[i].url < pageSlice[j].url
		}
	})
	fmt.Println("=============================")
  	fmt.Printf("REPORT for %v\n",baseURL)
	fmt.Println("=============================")
	for _, page := range pageSlice{
		fmt.Printf("Found %v internal links to %v\n",page.count,page.url)
	}

}



