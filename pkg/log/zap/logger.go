package zap

import (
	"github.com/consensysquorum/quorum-key-manager/pkg/log"
	"go.uber.org/zap"
)

type Logger struct {
	logger *zap.SugaredLogger
	cfg    *log.Config
}

func NewLogger(cfg *log.Config) (*Logger, error) {
	var logger *zap.Logger
	var err error

	if cfg.Mode == log.DevelopmentMode {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		return nil, err
	}

	return &Logger{logger: logger.Sugar(), cfg: cfg}, nil
}

func (l *Logger) Debug(msg string, keysAndValues ...interface{}) *Logger {
	l.logger.Debugw(msg, keysAndValues)
	return l
}

func (l *Logger) Info(msg string, keysAndValues ...interface{}) *Logger {
	l.logger.Infow(msg, keysAndValues)
	return l
}

func (l *Logger) Error(msg string, err error, keysAndValues ...interface{}) *Logger {
	l.logger.With("error", err).Errorw(msg, keysAndValues)
	return l
}

func (l *Logger) Panic(msg string, err error, keysAndValues ...interface{}) *Logger {
	l.logger.With("error", err).Panicw(msg, keysAndValues)
	return l
}

func (l *Logger) Fatal(msg string, err error, keysAndValues ...interface{}) *Logger {
	l.logger.With("error", err).Fatalw(msg, keysAndValues)
	return l
}
