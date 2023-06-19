package logger

import (
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func Init() {
	var err error
	l, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	logger = l.Sugar()
}

func Mock() {
	logger = zap.NewNop().Sugar()
}

func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}
