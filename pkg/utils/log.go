package utils

import (
	"fmt"
	"github.com/ghettovoice/gosip/log"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"os"
)

type MyLogger struct {
	Logger *log.LogrusLogger
	level  log.Level
}

func (ml *MyLogger) Level() string {
	switch ml.level {
	case log.PanicLevel:
		return "Panic"
	case log.FatalLevel:
		return "Fatal"
	case log.ErrorLevel:
		return "Error"
	case log.WarnLevel:
		return "Warn"
	case log.InfoLevel:
		return "Info"
	case log.DebugLevel:
		return "Debug"
	case log.TraceLevel:
		return "Trace"
	}
	return "Unkown"
}

var (
	loggers map[string]*MyLogger
)

func init() {
	loggers = make(map[string]*MyLogger)
}

func NewLogrusLogger(level log.Level, prefix string, fields log.Fields) log.Logger {
	if logger, found := loggers[prefix]; found {
		return logger.Logger.WithPrefix(prefix)
	}
	l := logrus.New()
	l.Level = logrus.ErrorLevel

	logFile, err := os.OpenFile(fmt.Sprintf("agent-%s.log", prefix), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		l.Fatalf("Failed to open log file: %s", err)
	}
	l.SetOutput(logFile)
	l.Formatter = &prefixed.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
		ForceColors:     true,
		ForceFormatting: true,
	}
	l.SetReportCaller(true)
	logger := log.NewLogrusLogger(l, "main", fields)
	loggers[prefix] = &MyLogger{
		Logger: logger,
		level:  level,
	}
	logger.SetLevel(uint32(level))
	return logger.WithPrefix(prefix)
}

func SetLogLevel(prefix string, level log.Level) error {
	if logger, found := loggers[prefix]; found {
		logger.level = level
		logger.Logger.SetLevel(uint32(level))
		return nil
	}
	return fmt.Errorf("logger [%v] not found", prefix)
}

func GetLoggers() map[string]*MyLogger {
	return loggers
}
