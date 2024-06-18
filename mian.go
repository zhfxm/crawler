package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main()  {
	url := "https://www.thepaper.cn/"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("http get err:%v", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("http error code: %v", resp.StatusCode)
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read io err:%v", err)
		return
	}
	numLinks := strings.Count(string(body), "<a")
	fmt.Printf("links count: %d \n", numLinks)
}