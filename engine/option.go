package engine

import (
	"github.com/zhfxm/simple-crawler/collect"
	"go.uber.org/zap"
)

type Option func(opts *options)

type options struct {
	WorkerCount int
	Logger      *zap.Logger
	Seeds       []*collect.Request
	Fetcher     collect.Fetcher
}

var defaultOption = options{
	WorkerCount: 3,
}

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.Logger = logger
	}
}

func WithFetch(f collect.Fetcher) Option {
	return func(opts *options) {
		opts.Fetcher = f
	}
}

func WithWorkCount(wc int) Option {
	return func(opts *options) {
		opts.WorkerCount = wc
	}
}

func WidthSeeds(seeds []*collect.Request) Option {
	return func(opts *options) {
		opts.Seeds = seeds
	}
}