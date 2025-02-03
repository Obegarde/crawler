package main
import(
	"fmt"
	"net/url"
)

func normalizeURL(inputURL string)(string,error){
	URLStruct, err := url.Parse(inputURL)
	if err != nil{
		fmt.Println(err)
		return "", err
	}
	normalizedPath := URLStruct.Path
	if string(normalizedPath[len(normalizedPath)-1])=="/"{
		normalizedPath = normalizedPath[:len(normalizedPath)-1]
	}
	normalizedURL := fmt.Sprintf("%v%v",URLStruct.Host,normalizedPath)
	fmt.Println(normalizedURL)
	return normalizedURL, nil
}
