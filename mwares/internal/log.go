package internal

import (
	"log"

	"go.uber.org/zap"
)

func InitZap() *zap.SugaredLogger {

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("init zap log failed: %v", err)
		return nil
	}

	defer func() {
		_ = logger.Sync()
	}()

	return logger.Sugar()
}
