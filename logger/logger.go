// Package logger is Heavily inspired by github.com/uber-go/zap
package logger

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	stdlog "log"
	"log/syslog"
	"net/url"
	"os"
	"time"
)

const (
	// LevelDebug info required while troubleshooting issues.
	LevelDebug Level = iota + 1
	// LevelInfo info required for future post-mortems.
	LevelInfo
)

// DefaultLevel is the level used when is not defined in configuration.
const DefaultLevel = LevelInfo

// Level represents the level type
type Level uint8

func (l Level) String() (string, error) {
	switch l {
	case LevelInfo:
		return "info", nil
	case LevelDebug:
		return "debug", nil
	}

	return "", errors.New("Invalid log level")
}

func stringToLevel(l string) (Level, error) {
	switch l {
	case "info":
		return LevelInfo, nil
	case "debug":
		return LevelDebug, nil
	}

	return 0, errors.New("Invalid log level")
}

// A Logger enables leveled, structured logging. All methods are safe for
// concurrent use.
type Logger interface {
	Info(string, ...Field)
	Debug(string, ...Field)
}

// Log implements Logger.
type log struct {
	marshaller     Marshaller
	writer         io.Writer
	thresholdLevel Level
}

// NewSyslogToLogstash creates a new rsyslog compatible logger.
func NewSyslogToLogstash(channel, application, environment string, l Level, facility syslog.Priority) (Logger, error) {
	var severity syslog.Priority
	switch l {
	case LevelDebug:
		severity = syslog.LOG_DEBUG
	case LevelInfo:
		severity = syslog.LOG_INFO
	}

	// The Priority value is calculated by first multiplying the Facility number by 8
	// and then adding the numerical value of the Severity.
	// See
	// 	https://tools.ietf.org/html/rfc3164#page-10
	// Or
	//	https://tools.ietf.org/html/rfc5424#page-11
	syslogWriter, err := syslog.New(facility|severity, fmt.Sprintf("%v.%v", application, environment))
	if err != nil {
		return nil, err
	}

	return newLog(
			newSyslogToLogstashMarshaller(channel, application, environment),
			syslogWriter,
			l,
		),
		nil
}

// NewDiscardAll creates a logger for testing purposes.
func NewDiscardAll() Logger {
	return newLog(
		newNOOPMarshaller(),
		ioutil.Discard,
		LevelInfo,
	)
}

// NewStdOut creates a logger which prints for the standard output.
func NewStdOut(channel, application, environment string, lvl Level) (Logger, error) {
	return newLog(
			newHumanMarshaller(channel, application, environment),
			os.Stdout,
			lvl,
		),
		nil
}

func newLog(marshaller Marshaller, writer io.Writer, thresholdLevel Level) Logger {
	return &log{
		marshaller,
		writer,
		thresholdLevel,
	}
}

// Info logs data with level Info.
func (l *log) Info(msg string, fields ...Field) {
	l.log(LevelInfo, msg, fields...)
}

// Debug data with level Debug.
func (l *log) Debug(msg string, fields ...Field) {
	l.log(LevelDebug, msg, fields...)
}

func (l *log) log(currentLevel Level, msg string, fields ...Field) {
	if currentLevel < l.thresholdLevel {
		return
	}

	entry := &Entry{
		message: msg,
		fields:  fields,
		level:   currentLevel,
		time:    time.Now(),
	}
	data, err := l.marshaller.Marshal(entry)

	if err == nil {
		_, err = l.writer.Write(data)
	}
	if err != nil {
		stdlog.Println(err)
	}
}

// NewLoggerFromDSN creates a new logger from a DSN configuration. If the configuration is not valid it panics.
func NewLoggerFromDSN(DSN, application, environment string) (Logger, error) {
	// param validation
	URL, err := url.Parse(DSN)
	if err != nil {
		return nil, err
	}

	lvl, err := levelFromURL(URL)
	if err != nil {
		return nil, err
	}

	var logger Logger
	channel := "Lernin"

	switch URL.Scheme {
	case "stdout":
		logger, err = NewStdOut(channel, application, environment, lvl)
	case "kibana":
		facility := syslog.LOG_LOCAL3 // We use LOCAL3 as default facility because of internal convention.
		logger, err = NewSyslogToLogstash(channel, application, environment, lvl, facility)
	case "discardall":
		logger = NewDiscardAll()
	default:
		return nil, errors.New("invalid logger publisher type")
	}

	if err != nil {
		return nil, err
	}

	return logger, nil
}

func levelFromURL(URL *url.URL) (Level, error) {
	lvl, exists := URL.Query()["level"]
	if !exists {
		return DefaultLevel, nil
	}

	return stringToLevel(lvl[0])
}
