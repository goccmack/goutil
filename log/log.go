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
