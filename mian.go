package main

import (
	"fmt"
	"log"
	"time"

	slog "github.com/zhfxm/crawler/log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	url := "www.google.com"
	fmt.Println("============logger.Info=======================")
	logger.Info("faild to fetch url",
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)
	fmt.Println("============sugar.Infow=======================")
	sugar := logger.Sugar()
	sugar.Infow("faild to fetch url",
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)

	fmt.Println("=============sugar1.Info======================")
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	logger1, _ := loggerConfig.Build()
	sugar1 := logger1.Sugar()
	sugar1.Info("logger config by self")

	fmt.Println("=============lumberjack======================")
	w := &lumberjack.Logger{
		Filename : "/Users/hui.zhang/Desktop/gplog.log",
		MaxSize: 500, // 日志最大大小，单位 M
		MaxBackups: 3, // 保留日志文件最大数量
		MaxAge: 28, // 保留旧日志文件最大天数
	}
	log.SetOutput(w)
	for i := 0; i < 2; i++ {
		log.Printf("This is log entry %d", i)
		time.Sleep(time.Second) // 模拟日志写入间隔
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), 
		zapcore.AddSync(w), 
		zap.InfoLevel,
	)
	logInstance := zap.New(core)
	zap.ReplaceGlobals(logInstance)
	for i := 0; i < 2; i++ {
		zap.L().Info("This is anothr log entry", zap.Int("entryNumber", i))
		time.Sleep(time.Second)
	}
	defer logInstance.Sync()

	fmt.Println("=============self log======================")
	plugin, c := slog.NewFilePlugin("/Users/hui.zhang/Desktop/slog.log", zap.InfoLevel)
	defer c.Close()
	slogger := slog.NewLogger(plugin)
	slogger.Info("slog info mesg")
}
