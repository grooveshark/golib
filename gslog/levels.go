package gslog

import (
	"os"
)

// logLevel is the severity of a log message
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

var levels = [...]string{
	"DEBUG",
	"INFO",
	"WARN",
	"ERROR",
	"FATAL",
}

// String returns the English name of the log level ("DEBUG", "INFO", ...)
func (ll LogLevel) String() string { return levels[ll] }

// Debug writes a message to the log with debug severity
func Debug(msg string, params ...interface{}) {
	log(DEBUG, &msg, params)
}

// Alias of Debug
var Debugf = Debug

// Info writes a message to the log with info severity
func Info(msg string, params ...interface{}) {
	log(INFO, &msg, params)
}

// Alias of Info
var Infof = Info

// Warn writes a message to the log with warn severity
func Warn(msg string, params ...interface{}) {
	log(WARN, &msg, params)
}

// Alias of Warn
var Warnf = Warn

// Error writes a message to the log with error serverity
func Error(msg string, params ...interface{}) {
	log(ERROR, &msg, params)
}

// Alias of Error
var Errorf = Error

// Fatal writes a message to the log with fatal severity, flushes any
// messages waiting to be written, and exits with a non-zero status
func Fatal(msg string, params ...interface{}) {
	log(FATAL, &msg, params)
	Flush()
	os.Exit(1)
}

// Alias of Fatal
var Fatalf = Fatal
