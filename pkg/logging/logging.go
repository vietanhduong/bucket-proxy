package logging

import (
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

var instance = defaultLogger()

type Formatter string

const (
	TextFormatter Formatter = "text"
	JSONFormatter Formatter = "json"
)

func WithField(key string, value any) *logrus.Entry {
	return instance.WithField(key, value)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return instance.WithFields(fields)
}

func SetLevel(level string) {
	var l logrus.Level
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		l = logrus.DebugLevel
	case "info":
		l = logrus.InfoLevel
	case "warn":
		l = logrus.WarnLevel
	case "error":
		l = logrus.ErrorLevel
	case "fatal":
		l = logrus.FatalLevel
	case "trace":
		l = logrus.TraceLevel
	default:
		l = logrus.InfoLevel
	}
	instance.SetLevel(l)
}

func SetOutput(output *os.File) {
	instance.SetOutput(output)
}

func Instance() *logrus.Logger {
	return instance
}

func SetFormatter(formatter Formatter, disableQuote bool) {
	format := formatter.ToLogrusFormatter()
	if f, ok := (format).(*logrus.TextFormatter); ok {
		f.DisableQuote = disableQuote
	}
	instance.SetFormatter(format)
}

func Debugf(format string, args ...interface{}) {
	instance.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	instance.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	instance.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	instance.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	instance.Fatalf(format, args...)
}

func Debug(args ...interface{}) {
	instance.Debug(args...)
}

func Info(args ...interface{}) {
	instance.Info(args...)
}

func Warn(args ...interface{}) {
	instance.Warn(args...)
}

func Error(args ...interface{}) {
	instance.Error(args...)
}

func Fatal(args ...interface{}) {
	instance.Fatal(args...)
}

func (f Formatter) ToLogrusFormatter() logrus.Formatter {
	switch f {
	case JSONFormatter:
		return &logrus.JSONFormatter{CallerPrettyfier: prettier}
	default:
		return &logrus.TextFormatter{
			DisableColors:    true,
			CallerPrettyfier: prettier,
		}
	}
}

var prettier = func(frame *runtime.Frame) (function string, file string) {
	fileName := path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
	return "", fileName
}

func defaultLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(TextFormatter.ToLogrusFormatter())
	logger.SetReportCaller(true)
	return logger
}
