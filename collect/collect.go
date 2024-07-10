package collect

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type Fetcher interface {
	Get(request *Request) ([]byte, error)
}

type BrowserFetch struct {
	Timeout time.Duration
}

func (b BrowserFetch) Get(request *Request) ([]byte, error) {
	client := &http.Client{
		Timeout: b.Timeout,
	}
	request.Log.Info("request url", zap.String("url", request.Url))
	req, err := http.NewRequest("GET", request.Url, nil)
	if err != nil {
		request.Log.Error("get url failed", zap.Error(err))
		return nil, fmt.Errorf("get url faild:%v", err)
	}

	if len(request.Cookie) > 0 {
		req.Header.Set("Cookie", request.Cookie)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	bodyReader := bufio.NewReader(resp.Body)
	e := DeterminEncoding(bodyReader, request.Log)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	return io.ReadAll(utf8Reader)
}

func DeterminEncoding(r *bufio.Reader, log *zap.Logger) encoding.Encoding {
	b, err := r.Peek(1024)
	if err != nil {
		log.Error("fetch error", zap.Error(err))
		return unicode.UTF8
	}
	e ,_ ,_ := charset.DetermineEncoding(b, "")
	return e
}