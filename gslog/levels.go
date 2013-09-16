package gslog

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

// Info writes a message to the log with info severity
func Info(msg string, params ...interface{}) {
	log(INFO, &msg, params)
}

// Warn writes a message to the log with warn severity
func Warn(msg string, params ...interface{}) {
	log(WARN, &msg, params)
}

// Error writes a message to the log with error serverity
func Error(msg string, params ...interface{}) {
	log(ERROR, &msg, params)
}

// Fatal writes a message to the log with fatal severity
func Fatal(msg string, params ...interface{}) {
	log(FATAL, &msg, params)
}
