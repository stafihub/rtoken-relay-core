package log

import "github.com/sirupsen/logrus"

type Logger interface {
	// Log a message at the given level with context key/value pairs
	Trace(msg string, ctx ...interface{})
	Debug(msg string, ctx ...interface{})
	Info(msg string, ctx ...interface{})
	Warn(msg string, ctx ...interface{})
	Error(msg string, ctx ...interface{})
}

type log struct{}

func NewLog(path string) Logger {
	if err := initLogFile(path); err != nil {
		panic(err)
	}
	return &log{}
}

func (l *log) Trace(msg string, ctx ...interface{}) {
	logrus.Trace(msg, ctx)
}

func (l *log) Debug(msg string, ctx ...interface{}) {
	logrus.Debug(msg, ctx)
}

func (l *log) Info(msg string, ctx ...interface{}) {
	logrus.Info(msg, ctx)
}

func (l *log) Warn(msg string, ctx ...interface{}) {
	logrus.Warn(msg, ctx)
}

func (l *log) Error(msg string, ctx ...interface{}) {
	logrus.Error(msg, ctx)
}
