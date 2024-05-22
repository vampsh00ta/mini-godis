package tests

import (
	"log"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type MockMethod struct {
	methodName string
	args       []any
	returns    []any
}

func LoadLoggerDev() *zap.SugaredLogger {
	loggerConfig := zap.NewDevelopmentConfig()

	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	loggerConfig.EncoderConfig.LevelKey = "CRITICAL"

	logger, err := loggerConfig.Build()
	if err != nil {
		log.Fatal(err)
	}
	sugar := logger.Sugar()
	return sugar
}
