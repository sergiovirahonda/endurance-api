package config

import (
	"os"

	"github.com/labstack/echo/v4"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	conf   = GetConfig()
	logger *zap.SugaredLogger
)

func initLog() {
	if logger != nil {
		return
	}
	encoderConfig := ecszap.NewDefaultEncoderConfig()
	core := ecszap.NewCore(encoderConfig, os.Stdout, getLevel(conf.Logger.Level))
	defaultLogger := zap.New(core, zap.AddCaller())
	logger = defaultLogger.Sugar()
}

func GetLogger() *zap.SugaredLogger {
	initLog()
	return logger
}

func GetLoggerFromContext(ctx echo.Context) *zap.SugaredLogger {
	requestID := ctx.Get("request_id")
	if requestID == nil {
		requestID = "no_request_id"
	}
	logger := ctx.Get("logger")
	if logger == nil {
		logger := GetLogger()
		logger.With("request_id", requestID)
		ctx.Set("logger", logger)
		return logger
	} else {
		logger := logger.(*zap.SugaredLogger)
		logger.With("request_id", requestID)
		ctx.Set("logger", logger)
	}
	return logger.(*zap.SugaredLogger)
}

func getLevel(lvl int64) zapcore.Level {
	switch lvl {
	case 0:
		return zap.InfoLevel
	case 1:
		return zap.WarnLevel
	case 2:
		return zap.ErrorLevel
	}
	return zap.DebugLevel
}
