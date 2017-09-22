package logger

import (
	"fmt"
	"time"
	"os"
	"strings"
)

var logLevel = LOG_INFO
var timeZone = time.Now().Format("Z07:00")

const (
	LOG_NONE  = 1 << iota
	LOG_ERROR = 1 << iota
	LOG_WARN  = 1 << iota
	LOG_INFO  = 1 << iota
	LOG_DEBUG = 1 << iota
)

const shorter_nano = "2006-01-02T15:04:05.9999"

func LevelAsString(level int) string {
	switch(level) {
	case LOG_NONE:
		return " NONE"
	case LOG_ERROR:
		return "ERROR"
	case LOG_WARN:
		return " WARN"
	case LOG_INFO:
		return " INFO"
	case LOG_DEBUG:
		return "DEBUG"
	default:
		return "UNDEF"
	}
}

func SetLevel(level int) {
	logLevel = level
}

func log(level int, format string, a []interface{}) (n int, err error) {
	if (level < logLevel) {
		return 0, nil
	}

	timeFormat := time.Now().Format(shorter_nano)
	for len(timeFormat) < len(shorter_nano) {
		timeFormat = timeFormat + "0"
	}

	return fmt.Printf("["+timeFormat+timeZone+"|telegraf-docker-sd|"+LevelAsString(level)+"] "+strings.TrimRight(format, "\n")+"\n", a...)
}

func Log(level int, format string, a ...interface{}) (n int, err error) {
	return log(level, format, a)
}

func Fatalf(format string, a ...interface{}) (n int, err error) {
	n, err = log(LOG_ERROR, format, a)
	os.Exit(255)
	return n, err
}

func Errorf(format string, a ...interface{}) (n int, err error) {
	return log(LOG_ERROR, format, a)
}

func Warnf(format string, a ...interface{}) (n int, err error) {
	return log(LOG_WARN, format, a)
}

func Infof(format string, a ...interface{}) (n int, err error) {
	return log(LOG_INFO, format, a)
}

func Debugf(format string, a ...interface{}) (n int, err error) {
	return log(LOG_DEBUG, format, a)
}
