package log

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"time"
)

// Priority of a logging message
type Priority int

const (
	// Exit is used to terminate with a specified UNIX exit code
	EXIT Priority = iota

	// PANIC is for irrecoverable program failure.
	// The logger calls os.Exit(1) after the message is logged and the logIF closed.
	PANIC

	// WARNING is for recoverable errors
	WARNING

	// INFO provides high-level information about program execution
	INFO

	// DEBUG is used for exeution tracing
	DEBUG
)

func (p Priority) String() string {
	switch p {
	case EXIT:
		return "EXIT"
	case PANIC:
		return "PANIC"
	case WARNING:
		return "WARNING"
	case INFO:
		return "INFO"
	case DEBUG:
		return "DEBUG"
	}
	panic(fmt.Sprintf("Invalid priority %d", p))
}

// ToPriority returns the Priority corresponding to str.
func ToPriority(str string) (Priority, error) {
	switch strings.ToUpper(str) {
	case "EXIT":
		return EXIT, nil
	case "PANIC":
		return PANIC, nil
	case "WARNING":
		return WARNING, nil
	case "INFO":
		return INFO, nil
	case "DEBUG":
		return DEBUG, nil
	}
	return DEBUG, errors.New("Invalid priority string " + str)
}

// Close is only necessary before os.Exit is called. Otherwise the logger will automatically
// close open files when the programe terminates. Calling log.Close() before the client program
// terminates will cause no harm.
func Close() {
	closeChan <- true

	// Give logger time to close files
	time.Sleep(time.Second)
}

// Exitf logs a formatted message followed by os.Exit(exitCode)
func Exitf(exitCode int, format string, a ...interface{}) {
	exitIF(exitCode, fmt.Sprintf(format, a...))
}

// Panicf logs a formatted message followed by a stack trace; flushes and closes the logIF file and
// then performs os.Exit(1)
func Panicf(format string, a ...interface{}) {
	panicIF(fmt.Sprintf(format, a...), getPanicStackTrace())
}

// Warningf logs a formatted message with priority Warning.
func Warningf(format string, a ...interface{}) {
	logIF(WARNING, format, a)
}

// Infof logs a formatted message with priority Info.
func Infof(format string, a ...interface{}) {
	logIF(INFO, format, a)
}

// Debugf logs a formatted message with priority Debug.
func Debugf(format string, a ...interface{}) {
	logIF(DEBUG, format, a)
}

// Exit logs a message followed by os.Exit(exitCode)
func Exit(exitCode int, msg string) {
	exitIF(exitCode, msg)
}

// Panic logs a message followed by a stack trace; flushes and closes the logIF file and
// then performs os.Exit(1)
func Panic(msg string) {
	panicIF(msg, getPanicStackTrace())
}

// Warning logs a message with priority Warning.
func Warning(msg string) {
	logIF(WARNING, msg, nil)
}

// Info logs a message with priority Info.
func Info(msg string) {
	logIF(INFO, msg, nil)
}

// Debug logs a message with priority Debug.
func Debug(msg string) {
	logIF(DEBUG, msg, nil)
}

// GetConfig returns the current logger configuration
func GetConfig() *Config {
	reply := make(chan *Config)
	getConfigChan <- reply
	select {
	case c := <-reply:
		return c
	case <-time.After(10 * time.Second):
		panic("Timeout waiting for log configuration")
	}
}

// SetConfig sets the configuration of the logger to priority, to use up to maxFiles files and to close
// files that exceed maxBytes
func SetConfig(maxFiles, maxBytes int, priority Priority) {
	setConfigChan <- &configMsg{
		maxFiles: maxFiles,
		maxBytes: maxBytes,
		priority: priority,
	}
}

// Suppress sets the list of files whose Debug messages are suppressed.Suppressed.
// If files is an empty string no files are suppressed.
// files is a comma separated list of file names.
// File names must not have a path.
// The ".go" extensions of the file names may be omitted.
//     E.g.: "file1,file2"
func Suppress(files string) {
	suppressChan <- files
}

func getPanicStackTrace() string {
	trc := make([]byte, 2048)
	length := runtime.Stack(trc, false)
	return string(trc[:length])
}
