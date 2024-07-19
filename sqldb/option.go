package sqldb

import (
	"go.uber.org/zap"
)

type Option func(opts *options)

type options struct {
	logger *zap.Logger
	sqlUrl string
}

var defaultOption = options{}

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func WidthSqlurl(sqlurl string) Option {
	return func(opts *options) {
		opts.sqlUrl = sqlurl
	}
}