package main

import(
	"net/http"
	"io"
	"fmt"
	"strings"
)

func getHTML(rawURL string)(string, error){
	res, err := http.Get(rawURL)
	if err != nil{
		return "",err
	}
	body,err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299{
		return "", fmt.Errorf("response failed with status code: %d and\nbody: %s", res.StatusCode, body)
	}
	if err != nil{
		return "", err
	}
	fmt.Println(res.Header["Content-Type"][0])
	if !strings.Contains(res.Header["Content-Type"][0], "text/html"){
		return "", err
	}	
	return string(body), nil
}
