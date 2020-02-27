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

The logger initialises and closes automatically.

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
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/goccmack/goutil/log/files"
)

/*** Interface to logger ***/

var (
	closeChan     = make(chan bool)
	exitChan      = make(chan *exitMsg)
	getConfigChan = make(chan chan *Config)
	logChan       = make(chan *logMsg, 1024)
	panicChan     = make(chan *panicMsg)
	setConfigChan = make(chan *configMsg)
	suppressChan  = make(chan string)
)

type configMsg struct {
	maxFiles int
	maxBytes int
	priority Priority
}

type exitMsg struct {
	file     string
	line     int
	exitCode int
	msg      string
}

type isSuppressedMsg struct {
	fileName string
	replyTo  chan bool
}

type logMsg struct {
	file     string
	line     int
	priority Priority
	format   string
	a        []interface{}
}

type panicMsg struct {
	file       string
	line       int
	msg        string
	stacktrace string
}

type logger struct {
	cfg *Config
	wtr *files.FileSet
}

func init() {
	go new(logger).run()
}

// exitIF is called from the logger interface routines
func exitIF(exitCode int, msg string) {
	file, line := getFileLine()
	exitChan <- &exitMsg{
		exitCode: exitCode,
		msg:      msg,
		file:     file,
		line:     line,
	}

	// wait for os.Exit()
	for {
		time.Sleep(time.Second)
	}
}

// logIF is called from the logger interface routines
func logIF(priority Priority, format string, a []interface{}) {
	lm := &logMsg{
		priority: priority,
		format:   format,
		a:        a,
	}
	lm.file, lm.line = getFileLine()

	defer recover()
	logChan <- lm
}

// panicIF is called from the logger interface routines
func panicIF(msg string, stackTrace string) {
	pm := &panicMsg{
		msg:        msg,
		stacktrace: stackTrace,
	}
	pm.file, pm.line = getFileLine()
	panicChan <- pm

	// wait for os.Exit(1)
	for {
		time.Sleep(time.Second)
	}
}

/*** logger class ***/

func (l *logger) close() {
	close(logChan)
	l.flushLogMsgs()
	l.wtr.Close()
}

func (l *logger) flushLogMsgs() {
	n := len(logChan)
	for i := 0; i < n; i++ {
		lm := <-logChan
		l.logMsg(lm.file, lm.line, lm.priority, lm.format, lm.a, "")
	}
}

func (l *logger) isSuppressed(file string, priority Priority) bool {
	if priority < DEBUG {
		return false
	}
	if l.cfg.SuppressedFiles == "" {
		return false
	}
	fns := strings.Split(file, ".")
	fn := strings.Join(fns[:len(fns)-1], ".")
	suppress := strings.Contains(l.cfg.SuppressedFiles, fn)
	return suppress
}

func (l *logger) logConfig() {
	fmt.Fprintf(l.wtr, "%s Log configuration:\n", time.Now().Format(time.RFC3339Nano))
	fmt.Fprintf(l.wtr, "  RootDir: %s\n", l.cfg.RootDir)
	fmt.Fprintf(l.wtr, "  NumFiles: %d\n", l.cfg.NumFiles)
	fmt.Fprintf(l.wtr, "  NumBytes: %d\n", l.cfg.FileNumBytes)
	fmt.Fprintf(l.wtr, "  Priority: %s\n", l.cfg.Priority)
	fmt.Fprintf(l.wtr, "  Suppress: %s\n", l.cfg.SuppressedFiles)
}

func (l *logger) logExit(file string, line int, exitCode int, msg string) {
	_, fname := path.Split(file)
	l.write(fmt.Sprintf("%s [EXIT %d] -%s, line %d- %s\n%s",
		time.Now().Format(time.RFC3339Nano),
		exitCode,
		fname, line,
		strings.TrimRight(msg, "\n"),
		""))
}

func (l *logger) logMsg(file string, line int, priority Priority,
	format string, a []interface{},
	stackTrace string) {

	_, fname := path.Split(file)
	if priority <= l.cfg.Priority && !l.isSuppressed(fname, priority) {
		msg := fmt.Sprintf(strings.TrimRight(format, "\n"), a...)
		l.write(fmt.Sprintf("%s [%s] -%s, line %d- %s\n%s",
			time.Now().Format(time.RFC3339Nano),
			priority,
			fname, line,
			msg,
			strings.TrimRight(stackTrace, "\n")))
	}
}

func (l *logger) run() {
	l.cfg = readConfigFile(true)
	l.wtr = files.New(l.cfg.RootDir, l.cfg.FileName, l.cfg.FileNumBytes, l.cfg.NumFiles)
	defer l.close()
	l.logConfig()

	refreshConfig := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-closeChan:
			return
		case msg := <-exitChan:
			l.logExit(msg.file, msg.line, msg.exitCode, msg.msg)
			l.close()
			os.Exit(msg.exitCode)
		case msg := <-logChan:
			l.logMsg(msg.file, msg.line, msg.priority, msg.format, msg.a, "")
		case msg := <-panicChan:
			l.logMsg(msg.file, msg.line, PANIC, msg.msg, nil, msg.stacktrace)
			l.close()
			os.Exit(1)
		case <-refreshConfig.C:
			newCfg := readConfigFile(false)
			if !l.cfg.Equal(newCfg) {
				l.cfg = newCfg
				l.wtr.SetConfig(l.cfg.NumFiles, l.cfg.FileNumBytes)
				l.logConfig()
			}
		case cm := <-setConfigChan:
			l.cfg.NumFiles = cm.maxFiles
			l.cfg.FileNumBytes = cm.maxBytes
			l.cfg.Priority = cm.priority
			l.flushLogMsgs()
			l.wtr.SetConfig(cm.maxFiles, cm.maxBytes)
			l.logConfig()
		case replyTo := <-getConfigChan:
			replyTo <- l.cfg.Clone()
		case s := <-suppressChan:
			l.cfg.SuppressedFiles = s
			l.flushLogMsgs()
			l.logConfig()
		}
	}
}

func (l *logger) write(msg string) {
	if _, err := l.wtr.Write(([]byte)(msg)); err != nil {
		panic(err)
	}
}

/***** Utility ******/

func getFileLine() (file string, line int) {
	_, file, line, _ = runtime.Caller(3)
	return file, line
}
