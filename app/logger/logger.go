package logger

import (
	"fmt"
	"time"
	"os"
	"strings"
)

var maxLevel = LOG_INFO
var timeZone = time.Now().Format("Z07:00")

const (
	LOG_NONE  = 0
	LOG_ERROR = iota
	LOG_WARN  = iota
	LOG_INFO  = iota
	LOG_DEBUG = iota
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

func LevelFromString(level string) int {
	switch(level) {
	case "none":
		return LOG_NONE
	case "error":
		return LOG_ERROR
	case "warn":
		return LOG_WARN
	case "info":
		return LOG_INFO
	case "debug":
		return LOG_DEBUG
	default:
		return 0
	}
}

func SetLevel(level int) {
	maxLevel = level
}

func log(level int, format string, a []interface{}) (n int, err error) {
	if level > maxLevel {
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
