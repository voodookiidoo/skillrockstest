package logger

import "go.uber.org/zap"

func DefaultLogger() *zap.Logger {
	return zap.Must(zap.NewProduction())
}
