package main
import(
	"fmt"
	"net/url"
	"strings"
)

func normalizeURL(inputURL string)(string,error){
	URLStruct, err := url.Parse(inputURL)
	if err != nil{
		return "", fmt.Errorf("couldn't parse URL:%v",err)
	}
	normalizedPath := URLStruct.Path
	if string(normalizedPath[len(normalizedPath)-1])=="/"{
		normalizedPath = normalizedPath[:len(normalizedPath)-1]
	}
	normalizedURL := strings.ToLower(fmt.Sprintf("%v%v",URLStruct.Host,normalizedPath))	
	return normalizedURL, nil
}
