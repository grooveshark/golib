package gslog

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type message string

type logger struct {
	sync.RWMutex
	fileHandle *os.File
	filePath	string
	messages   chan *message
	flushCh    chan chan bool
	minLevel   LogLevel
}

var l *logger = nil

// init initializes the logger
func init() {
	l = &logger{}
	l.fileHandle = os.Stderr
	l.messages = make(chan *message, 1024)
	l.flushCh = make(chan chan bool)

	go handleMessages()
}

// SetMinimumLevel sets the minimum log level that will be output to the error
// log.  The capitilization of level does not matter.  Any log message of a
// lower severity will be silently ignored.  Default is DEBUG.
func SetMinimumLevel(level string) error {
	l.RWMutex.Lock()
	defer l.RWMutex.Unlock()
	var levelConstant LogLevel
	switch strings.ToUpper(level) {
	case "DEBUG":
		levelConstant = DEBUG
	case "INFO":
		levelConstant = INFO
	case "WARN":
		levelConstant = WARN
	case "ERROR":
		levelConstant = ERROR
	case "FATAL":
		levelConstant = FATAL
	default:
		return errors.New("invalid log level")
	}
	l.minLevel = levelConstant
	return nil
}

// SetLogFile sets the file to which messages will be logged to.  Default is STDOUT.
func SetLogFile(path string) error {
	l.RWMutex.Lock()
	defer l.RWMutex.Unlock()

	var fh *os.File
	var err error = nil

	if path == "stderr" {
		// Clear filePath so it doesn't get reopenned
		l.filePath = ""
		if l.fileHandle == os.Stderr {
			return nil
		}
		fh = os.Stderr
	} else {
		l.filePath = path
		flags := os.O_APPEND | os.O_CREATE | os.O_WRONLY
		fh, err = os.OpenFile(path, flags, 0644)
	}

	if err != nil {
		return err
	} else {
		l.fileHandle.Sync()
		l.fileHandle.Close()
		l.fileHandle = fh
	}

	return nil
}

// Flushes the log of at least 100 milliseconds worth of entries
func Flush() {
	retCh := make(chan bool)
	l.flushCh <- retCh
	<-retCh
}

// logString generates the string which will be written to the logfile
func logString(level LogLevel, msg *string, params []interface{}) string {
	ts := time.Now().Format(time.RFC1123)
	msgf := fmt.Sprintf(*msg, params...)
	return fmt.Sprintf("[%s] %s --> %s\n", ts, level.String(), msgf)
}

// log writes a message to the log file with the provided severity
func log(level LogLevel, msg *string, params []interface{}) {
	l.RWMutex.RLock()
	defer l.RWMutex.RUnlock()

	if level < l.minLevel {
		return
	}

	msgobj := message(logString(level, msg, params))
	l.messages <- &msgobj
}

// writeMessage writes the given message to the currently active logfile, or if
// there is an error doing that writes it to stderr along with whatever error
// writing to the logfile gave
func writeMessage(msg *message) {
	l.RWMutex.RLock()
	defer l.RWMutex.RUnlock()

	_, err := l.fileHandle.WriteString(string(*msg))
	if err != nil { // Try to reopen and write one more time
		// Assumes filePath is set as it should be impossible to error writing
		// to stderr
		SetLogFile(l.filePath)
		_, err = l.fileHandle.WriteString(string(*msg))
	}
	if err != nil { // If still an error writing, set to stderr
		errstr := err.Error()
		str := logString(ERROR, &errstr, []interface{}{})
		os.Stderr.WriteString(str)
		os.Stderr.WriteString(string(*msg))
	}
}

// handleMessages reads messages from a channel and writes them to the log file
// we have open. If no new log messages have come in for 100 milliseconds then
// we see if anyone has made a flush request. If they have then we send them
// back a true to indicate that we haven't had any messages for 100
// milliseconds, then we go again.
func handleMessages() {
	for {
		select {

		case msg, ok := <-l.messages:
			if !ok {
				panic("logger messages channel was closed!")
			} else {
				writeMessage(msg)
			}

		case <-time.After(100 * time.Millisecond):
			select {
			case flushret := <-l.flushCh:
				flushret <- true
			default: // Oh well
			}

		}
	}
}
