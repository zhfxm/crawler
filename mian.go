package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main()  {
	proxy, err := NewProxy()
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", ProxyRequestHandly(proxy))
	log.Fatal(http.ListenAndServe(":8181", nil))
}

func NewProxy() (*httputil.ReverseProxy, error) {
	targetHost := "http://localhost:8080"
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(url)
	return proxy, err
}

func ProxyRequestHandly(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}