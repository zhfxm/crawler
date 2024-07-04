package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func DefaultEncoderConfig() zapcore.EncoderConfig {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return encoderConfig
}

func DefaultEncode() zapcore.Encoder {
	return zapcore.NewJSONEncoder(DefaultEncoderConfig())
}

func DefaultOption() []zap.Option {
	var stackTraceLevel zap.LevelEnablerFunc = func(level zapcore.Level) bool {
		return level >= zapcore.DPanicLevel
	}
	return []zap.Option{
		zap.AddCaller(),
		zap.AddStacktrace(stackTraceLevel),
	}
}

func DefaultLumberjackLogger() *lumberjack.Logger {
	return & lumberjack.Logger{
		MaxSize: 300,
		LocalTime: true,
		Compress: true,
	}
}