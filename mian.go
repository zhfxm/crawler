package main

import (
	"fmt"

	"github.com/zhfxm/crawler/collect"
)

func main()  {
	url := "https://www.thepaper.cn/"
	var f = collect.BaseFetch{}
	b, err := f.Get(url)
	if err != nil {
		fmt.Printf("collect get error: %v\n", err)
	}
	fmt.Printf("get info: %s \n", string(b))
}