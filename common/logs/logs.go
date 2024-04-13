package logs

import (
	"common/config"
	"github.com/charmbracelet/log"
	"os"
	"time"
)

var logger *log.Logger

func InitLog(appName string) {
	logger = log.New(os.Stderr)
	if config.Conf.Log.Level == "DEBUG" {
		logger.SetLevel(log.DebugLevel)
	} else {
		logger.SetLevel(log.InfoLevel)
	}
	logger.SetPrefix(appName)
	logger.SetReportTimestamp(true)
	logger.SetTimeFormat(time.DateTime)
}

func Fatal(format string, v ...interface{}) {
	if len(v) == 0 {
		logger.Fatal(format)
	} else {
		logger.Fatalf(format, v...)
	}
}

func Info(format string, v ...interface{}) {
	if len(v) == 0 {
		logger.Info(format)
	} else {
		logger.Infof(format, v...)
	}
}

func Warning(format string, v ...interface{}) {
	if len(v) == 0 {
		logger.Warn(format)
	} else {
		logger.Warnf(format, v...)
	}
}

func Debug(format string, v ...interface{}) {
	if len(v) == 0 {
		logger.Debug(format)
	} else {
		logger.Debugf(format, v...)
	}
}

func Error(format string, v ...interface{}) {
	if len(v) == 0 {
		logger.Error(format)
	} else {
		logger.Errorf(format, v...)
	}
}
