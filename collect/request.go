package collect

import (
	"crypto/md5"
	"encoding/hex"
)

type Request struct {
	Url       string
	Method    string
	Cookie    string
	Priority  int
	ParseFunc func([]byte, *Request) ParseResult
}

type ParseResult struct {
	Requests []*Request
	Items    []interface{}
}

// 请求唯一标识码
func (r *Request) Unique() string {
	block := md5.Sum([]byte(r.Url + r.Method))
	return hex.EncodeToString(block[:])
}
