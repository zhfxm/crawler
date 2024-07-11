package main

import (
	"fmt"
	"time"

	"github.com/zhfxm/simple-crawler/collect"
	"github.com/zhfxm/simple-crawler/engine"
	"github.com/zhfxm/simple-crawler/log"
	"github.com/zhfxm/simple-crawler/parse/doubangroup"
	"go.uber.org/zap/zapcore"
)

func main() {

	plugin := log.NewStdoutPlugin(zapcore.InfoLevel)
	logger := log.NewLogger(plugin)
	logger.Info("logger init end")

	cookie := "viewed=\"1007305\"; bid=QGFx5rUryE0; _pk_id.100001.8cb4=c8407de68c359e04.1720513239.; __utmc=30149280; dbcl2=\"270053320:gi+WD2F5+Ic\"; ck=LlMn; push_noty_num=0; push_doumail_num=0; __utmz=30149280.1720576818.3.2.utmcsr=accounts.douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; __utmv=30149280.27005; douban-fav-remind=1; _pk_ref.100001.8cb4=%5B%22%22%2C%22%22%2C1720667364%2C%22https%3A%2F%2Faccounts.douban.com%2F%22%5D; _pk_ses.100001.8cb4=1; __yadk_uid=TUF1KmzWhyGajolSzhk5D6wFjHxGnxoU; __utma=30149280.471027257.1718694576.1720600444.1720667366.5; __utmt=1; __utmb=30149280.7.5.1720667366"

	var seeds []*collect.Request
	for i := 0; i < 100; i += 25 {
		str := fmt.Sprintf("https://www.douban.com/group/szsh/discussion?start=%d&type=new", i)
		seeds = append(seeds, &collect.Request{
			Url:       str,
			ParseFunc: doubangroup.ParseURL,
			Cookie:    cookie,
		})
	}
	logger.Sugar().Info("seeds size:", len(seeds))

	var f collect.Fetcher = collect.BrowserFetch{
		Timeout: 30 * time.Second,
		Logger:  logger,
	}

	e := engine.NewSchedule(
		engine.WidthSeeds(seeds),
		engine.WithFetch(f),
		engine.WithLogger(logger),
		engine.WithWorkCount(3),
	)
	e.Run()
}
