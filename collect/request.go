package collect

import "go.uber.org/zap"

type Request struct {
	Url       string
	Cookie    string
	ParseFunc func([]byte, *Request) ParseResult
	Log       *zap.Logger
}

type ParseResult struct {
	Requests []*Request
	Items    []interface{}
}
