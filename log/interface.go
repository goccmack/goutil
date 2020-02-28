//  Copyright 2020 Marius Ackerman
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

/*
Package log supports logging to a managed set of log files.

The logger is configured by a JSON file called <component>.log.config. <component> is the name of
the go binary executable (os.Executable()).
The logger looks for the log config file in $PWD.

The following is a JSON configuration structure containing the default logger paramerters:

	{
		"RootDir": "/usr/local/var/log",
		"NumFiles": 3,
		"FileNumBytes": 1000000,
		"Priority": "INFO",
		"SuppressedFiles": ""
	}

If the working directory does not contain a log.config file the logger uses these parameters. All
fields of log.config are optional. The logger will used default values for missing parameters.

The logger reads log.config periodically. The logger uses changed parameters. The log.config can
be changed while the program is running and further logging reflects the changed log.config.

The logger initialises and closes automatically but log.Close() should be called to ensure that
the last logged items are properly flushed before the program terminates.

log.Panic(...) and log.Panicf(...) log a stack trace in addition to the log message, close the
open log file and call os.Exit(1).

The logger will automatically tag every log message with time, the source file and
line number of the call to log.

If the logger priority is DEBUG the logger can be configured to suppress the debug messages from one
or more files by providing a comma separated string to the suppressFilesDebug parameter
of log.Init(...). The components of the string correspond to file names without extension. E.g.:
"pkga,pkgb" will suppress debug messages from pkga.go and pkgb.go.

Package log supports four Priority levels in decreasing order of priority:
Panic, Warning, Info, Debug. The logger instance has two methods to log a message of each Priority:
<Priority> and <Priority>f, e.g.: Info and Infof. <Priority> takes a string parameter, while
<Priority>f takes a format string followed by a list of parameters. The format of the <Priorty>f
format string parameter is the same as for fmt.Printf.

The logger is initialised with the minimum priority level. The logger discards all messages logged with a
lower priority.

The directory log/examples demonstrates some featurs of log:

	examples/basic:
		Shows the basic usage of the log package.

	examples/set_config:
		Shows a change of log file size and debug level on a live program.

	examples/suppress_files:
		Shows the suppression of debug message from a selection of files

examples/basic:

	log.config:
	{
		"RootDir": "logs",
	}

Rootdir sets the logging directory to the directory "logs" under the working directory.
The other Config parameters will take their default values.

	package main

	import (
		"github.com/goccmack/goutil/log/example/pkga"

		"github.com/goccmack/goutil/log"
	)

	package main

	import (
		"github.com/goccmack/goutil/log/examples/basic/pkga"

		"github.com/goccmack/goutil/log"
	)

	func main() {
		// Log something
		log.Info("This message WILL appear in the log")

		pkga.Go()
		// Give pkga time to log a message
		time.Sleep(time.Millisecond)

		log.Debug("This message will NOT appear in the log")
		log.Panic("This is a panic")
	}

	======================

	package pkga

	import (
		"github.com/goccmack/goutil/log"
	)

	func Go() {
		// Log something
		log.Debug("This is pkga")
	}

The resulting log is:

	File set configuration @ 2018-12-20T12:07:58.555278+01:00
	Maximum file size 1000000 bytes
	Maximum 3 files
	2018-12-20T12:07:58.555087+01:00 Log configuration:
	RootDir: logs
	NumFiles: 3
	NumBytes: 1000000
	Priority: INFO
	Suppress:
	2018-12-20T12:07:58.555464+01:00 [INFO] -main.go, line 14- This message WILL appear in the log
	2018-12-20T12:07:58.555482+01:00 [INFO] -pgka.go, line 9- This is pkga
	2018-12-20T12:07:58.555495+01:00 [PANIC] -main.go, line 19- This is a panic
	goroutine 1 [running]:
	github.com/goccmack/goutil/log.getPanicStackTrace(0x10bea8e, 0x3)
		.../github.com/goccmack/goutil/log/interface.go:138 +0x74
	github.com/goccmack/goutil/log.Panic(0x10fac35, 0xf)
		.../github.com/goccmack/goutil/log/interface.go:96 +0x22
	main.main()
		.../github.com/goccmack/goutil/log/examples/basic/main.go:19 +0x8d
*/
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
