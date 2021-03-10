package log

import (
	"github.com/sirupsen/logrus"
)

type LoggerLevel string

const (
	InfoLevel  LoggerLevel = "info"
	ErrorLevel LoggerLevel = "error"
	DebugLevel LoggerLevel = "debug"
	WarnLevel  LoggerLevel = "warn"
	TraceLevel LoggerLevel = "trace"
)

type Config struct {
	Level     LoggerLevel
	Timestamp bool
}

func (l LoggerLevel) logrusLvl() logrus.Level {
	switch l {
	case ErrorLevel:
		return logrus.ErrorLevel
	case WarnLevel:
		return logrus.WarnLevel
	case InfoLevel:
		return logrus.InfoLevel
	case DebugLevel:
		return logrus.DebugLevel
	case TraceLevel:
		return logrus.TraceLevel
	}

	return logrus.InfoLevel
}
