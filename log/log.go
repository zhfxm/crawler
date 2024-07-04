package log

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Plugin = zapcore.Core

func NewLogger(plugin zapcore.Core, options ...zap.Option) *zap.Logger {
	return zap.New(plugin, append(DefaultOption(), options...)...)
}

func NewPlugin(writer zapcore.WriteSyncer, enabler zapcore.LevelEnabler) Plugin {
	std := zapcore.AddSync(os.Stdout)
	mw := zapcore.NewMultiWriteSyncer(writer, std)
	return zapcore.NewCore(DefaultEncode(), mw, enabler)
}

func NewStdoutPlugin(enabler zapcore.LevelEnabler) Plugin {
	return NewPlugin(zapcore.Lock(zapcore.AddSync(os.Stdout)), enabler)
}

func NewStderrPlugin(enabler zapcore.LevelEnabler) Plugin {
	return NewPlugin(zapcore.Lock(zapcore.AddSync(os.Stderr)), enabler)
}

func NewFilePlugin(filepath string, enabler zapcore.LevelEnabler) (Plugin, io.Closer) {
	write := DefaultLumberjackLogger()
	write.Filename = filepath
	return NewPlugin(zapcore.AddSync(write), enabler), write
}