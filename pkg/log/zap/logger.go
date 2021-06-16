package zap

import (
	"github.com/consensysquorum/quorum-key-manager/pkg/log"
	"go.uber.org/zap"
)

type Logger struct {
	logger *zap.SugaredLogger
	cfg    *log.Config
}

var _ log.Logger = &Logger{}

func NewLogger(cfg *log.Config) (*Logger, error) {
	var logger *zap.Logger
	var err error

	switch cfg.Mode {
	case log.DevelopmentMode:
		logger, err = zap.NewDevelopment()
	case log.ProductionMode:
		logger, err = zap.NewProduction()
	default:
		logger, err = zap.NewProduction()
	}
	if err != nil {
		return nil, err
	}

	return &Logger{logger: logger.Sugar(), cfg: cfg}, nil
}

func (l *Logger) Debug(msg string, keysAndValues ...interface{}) log.Logger {
	l.logger.Debugw(msg, keysAndValues)
	return l
}

func (l *Logger) Info(msg string, keysAndValues ...interface{}) log.Logger {
	l.logger.Infow(msg, keysAndValues)
	return l
}

func (l *Logger) Warn(msg string, keysAndValues ...interface{}) log.Logger {
	l.logger.Warnw(msg, keysAndValues)
	return l
}

func (l *Logger) Error(msg string, keysAndValues ...interface{}) log.Logger {
	l.logger.Errorw(msg, keysAndValues)
	return l
}

func (l *Logger) Panic(msg string, keysAndValues ...interface{}) log.Logger {
	l.logger.Panicw(msg, keysAndValues)
	return l
}

func (l *Logger) Fatal(msg string, keysAndValues ...interface{}) log.Logger {
	l.logger.Fatalw(msg, keysAndValues)
	return l
}

func (l Logger) WithError(err error) log.Logger {
	l.logger = l.logger.With("error", err)
	return &l
}

func (l Logger) With(args ...interface{}) log.Logger {
	l.logger = l.logger.With(args)
	return &l
}

func (l Logger) WithComponent(component string) log.Logger {
	l.logger.Desugar().Named(component).Sugar()
	return &l
}

func (l *Logger) Write(p []byte) (n int, err error) {
	switch l.cfg.Level {
	case log.DebugLevel:
		l.Debug(string(p))
	case log.InfoLevel:
		l.Info(string(p))
	case log.WarnLevel:
		l.Warn(string(p))
	case log.ErrorLevel:
		l.Error(string(p))
	case log.PanicLevel:
		l.Panic(string(p))
	default:
		l.Info(string(p))
	}

	return 0, nil
}

func (l Logger) Sync() error {
	return l.logger.Sync()
}
