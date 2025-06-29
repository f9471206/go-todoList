package utils

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger(isProduction bool) {
	var err error
	if isProduction {
		logDir := "logs"
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			panic(err)
		}

		logFile := filepath.Join(logDir, "app.log")
		file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}

		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(file),
			zap.InfoLevel,
		)

		Logger = zap.New(core)
	} else {
		Logger, err = zap.NewDevelopment()
	}
	if err != nil {
		panic(err)
	}
}
