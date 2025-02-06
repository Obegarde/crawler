package main
import(
	"fmt"
	"os"
)
func main(){
	 argsWithProg := os.Args
	if len(argsWithProg) < 4{
		fmt.Println("Must provide both website and max pages to crawl")
		os.Exit(1)
	}else if len(argsWithProg) > 4{
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}else{
		cfg, err := configure(argsWithProg[1],argsWithProg[2],argsWithProg[3])
		if err != nil{
			fmt.Println("Failed to configure crawler")
			os.Exit(1)
		}
		fmt.Println("starting crawl")
		cfg.crawlPage(argsWithProg[1])
		cfg.wg.Wait()
		fmt.Println(cfg.pages)
		return
	}

}
