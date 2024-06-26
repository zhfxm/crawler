package collect

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
)

type Fetcher interface {
	Get(url string) ([]byte, error)
}

type BaseFetch struct {

}

func (BaseFetch) Get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error status code: %d", resp.StatusCode)
	}
	bodyReader := bufio.NewReader(resp.Body)
	return io.ReadAll(bodyReader)
}