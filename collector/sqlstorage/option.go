package sqlstorage

import "go.uber.org/zap"

type Option func(opt *options)

type options struct {
	logger     *zap.Logger
	sqlurl     string
	BatchCount int
}

var defaultOption = options{
	BatchCount: 1000,
}

func WithLogger(logger *zap.Logger) Option {
	return func(opt *options) {
		opt.logger = logger
	}
}

func WidthSqlurl(sqlurl string) Option {
	return func(opt *options) {
		opt.sqlurl = sqlurl
	}
}

func WidthBatchCount(batchCount int) Option {
	return func(opt *options) {
		opt.BatchCount = batchCount
	}
}
