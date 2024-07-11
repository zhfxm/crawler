package doubangroup

import (
	"regexp"

	"github.com/zhfxm/simple-crawler/collect"
)

const cityListRe = `(https://www.douban.com/group/topic/[0-9a-z]+/)"[^>]*>([^<]+)</a>`

func ParseURL(contents []byte, req *collect.Request) collect.ParseResult {
	re := regexp.MustCompile(cityListRe)
	matchs := re.FindAllSubmatch(contents, -1)
	result := collect.ParseResult{}
	for _, m := range matchs {
		u := string(m[1])
		result.Requests = append(result.Requests, 
			&collect.Request{
				Url: u,
				Cookie: req.Cookie,
				ParseFunc: func(c []byte, request *collect.Request) collect.ParseResult {
					return GetContent(c, u)
				},
			})
	}
	return result
}

const ContentRe = `<div class="topic-content">[\s\S]*?阳台[\s\S]*?<div`

func GetContent(contents []byte, url string) collect.ParseResult {
	re := regexp.MustCompile(ContentRe)
	ok := re.Match(contents)
	if !ok {
		return collect.ParseResult{
			Items: []interface{}{},
		}
	}
	result := collect.ParseResult{
		Items: []interface{}{url},
	}
	return result
}