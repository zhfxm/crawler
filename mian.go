package main

import (
	"time"

	"github.com/zhfxm/simple-crawler/collect"
	"github.com/zhfxm/simple-crawler/collector/sqlstorage"
	"github.com/zhfxm/simple-crawler/engine"
	"github.com/zhfxm/simple-crawler/log"
	"github.com/zhfxm/simple-crawler/parse/pool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {

	plugin := log.NewStdoutPlugin(zapcore.InfoLevel)
	logger := log.NewLogger(plugin)
	logger.Info("logger init end")

	dbUrl := "root:root@tcp(localhost:3306)/cloud_crawler?charset=utf8"
	storage, err := sqlstorage.NewStorage(
		sqlstorage.WidthBatchCount(1),
		sqlstorage.WidthSqlurl(dbUrl),
		sqlstorage.WithLogger(logger),
	)
	if err != nil {
		logger.Error("sql storage init error", zap.Error(err))
	}

	var seeds = make([]*collect.Request, 0, 1000)
	for i := 0; i < 1; i++ {
		// str := fmt.Sprintf("https://www.douban.com/group/szsh/discussion?start=%d&type=new", i)
		str := "https://www.antpool.com/auth/v3/observer/api/hash/query?accessKey=SM5yTbAwBFk32orAZRKW&coinType=BTC&observerUserId=JamjeeH2"
		seeds = append(seeds, &collect.Request{
			Url:       str,
			ParseFunc: pool.AntpoolCurrentParseFunc,
			Logger:    logger,
		})
	}
	logger.Sugar().Info("seeds size:", len(seeds))

	var f collect.Fetcher = collect.BrowserFetch{
		Timeout: 30 * time.Second,
		Logger:  logger,
	}

	s := engine.NewSchedule(logger)
	e := engine.NewEngine(
		engine.WidthScheduler(s),
		engine.WidthSeeds(seeds),
		engine.WithFetch(f),
		engine.WithLogger(logger),
		engine.WithWorkCount(10),
		engine.WidthStorage(storage),
	)
	e.Run()
}
