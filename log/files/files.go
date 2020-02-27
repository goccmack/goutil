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
Package files implements a managed fileset writer/closer.
*/
package files

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type FileSet struct {
	closeChan       chan chan bool
	currentFile     *os.File
	currentFileSize int
	logDir          string
	logName         string
	maxFileSize     int
	maxNumFiles     int
	msgChan         chan *writeRequest
	setConfigChan   chan *setConfig
}

type setConfig struct {
	fileSize int
	numFiles int
	replyTo  chan bool
}

type writeRequest struct {
	msg   []byte
	reply chan *writeResponse
}

type writeResponse struct {
	n   int
	err error
}

func New(logDir, logName string, maxFileSize, maxNumFiles int) *FileSet {
	fmt.Fprintf(os.Stdout, "Log directory: %s\n", logDir)
	fs := &FileSet{
		closeChan:     make(chan chan bool, 1),
		logDir:        logDir,
		logName:       logName,
		maxFileSize:   maxFileSize,
		maxNumFiles:   maxNumFiles,
		msgChan:       make(chan *writeRequest, 1024),
		setConfigChan: make(chan *setConfig),
	}
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		panic(err)
	}
	go fs.run()
	return fs
}

func (fs *FileSet) Close() {
	reply := make(chan bool)
	fs.closeChan <- reply
	select {
	case <-reply:
		return
	case <-time.After(time.Second):
		panic("timeout")
	}
}

// ListLogFiles returns the logfiles of logname in logDir sorted from oldest to newest
func ListLogFiles(logDir, logName string) []string {
	froot := filepath.Join(logDir, logName)
	pattern := fmt.Sprintf("%s*.log", froot)
	fs, err := filepath.Glob(pattern)
	if err != nil {
		panic(err)
	}
	sort.Strings(fs)

	return fs
}

// SetConfig sets the maximum number of log files to numfiles and
// the maximum file size to filesize bytes.
func (fs *FileSet) SetConfig(numFiles, fileSize int) {
	reply := make(chan bool)
	fs.setConfigChan <- &setConfig{
		numFiles: numFiles,
		fileSize: fileSize,
		replyTo:  reply,
	}
	select {
	case <-reply:
	case <-time.After(time.Second):
		panic("Timeout waiting for files to set config")
	}
}

func (fs *FileSet) Write(buf []byte) (int, error) {
	reply := make(chan *writeResponse)
	fs.msgChan <- &writeRequest{
		msg:   buf,
		reply: reply,
	}
	select {
	case <-time.After(time.Second):
		panic("Timeout")
	case rep := <-reply:
		return rep.n, rep.err
	}
}

/*** FileSet ***/

func (fs *FileSet) close() {
	close(fs.msgChan)
	for msg := range fs.msgChan {
		fs.log(msg.msg)
	}

	fname := fs.currentFile.Name()
	if fs.currentFileSize < 1 {
		fs.rmFile(fname)
	}

	fs.currentFile.Close()
}

func (fs *FileSet) listLogFiles() []string {
	return ListLogFiles(fs.logDir, fs.logName)
}

func (fs *FileSet) log(buf []byte) *writeResponse {
	n, err := fs.currentFile.Write(buf)
	if err == nil {
		fs.currentFileSize += len(buf)
		if fs.currentFileSize >= fs.maxFileSize {
			fs.rotate()
		}
	}
	return &writeResponse{n, err}
}

func (fs *FileSet) logConfig() {
	fmt.Fprintf(fs.currentFile, "File set configuration @ %s\n", time.Now().Format(time.RFC3339Nano))
	fmt.Fprintf(fs.currentFile, "Maximum file size %d bytes\n", fs.maxFileSize)
	fmt.Fprintf(fs.currentFile, "Maximum %d files\n", fs.maxNumFiles)
}

func (fs *FileSet) newFile() {
	tm := time.Now().Format(time.RFC3339Nano)
	fname := filepath.Join(fs.logDir,
		fmt.Sprintf("%s_%s.log", fs.logName, tm))

	var err error
	// fs.currentFile, err = os.Create(filepath.Join(fs.logDir, fname))
	fs.currentFile, err = os.Create(fname)
	if err != nil {
		panic(err)
	}

	fs.logConfig()
}

func (fs *FileSet) rmFile(fname string) {
	if err := os.Remove(fname); err != nil {
		panic(err)
	}
}

func (fs *FileSet) rotate() {
	if fs.currentFile != nil {
		fs.currentFile.Close()
	}
	logFiles := fs.listLogFiles()
	delete := len(logFiles) - fs.maxNumFiles + 1
	for i := 0; i < delete; i++ {
		fs.rmFile(logFiles[i])
	}
	fs.newFile()
	fs.currentFileSize = 0
}

func (fs *FileSet) run() {
	fs.rotate()
	for {
		select {
		case done := <-fs.closeChan:
			fs.close()
			done <- true
			return
		case cfg := <-fs.setConfigChan:
			fs.setConfig(cfg)
			cfg.replyTo <- true
		case msg := <-fs.msgChan:
			msg.reply <- fs.log(msg.msg)
		}
	}
}

func (fs *FileSet) setConfig(cfg *setConfig) {
	fs.maxFileSize = cfg.fileSize
	fs.maxNumFiles = cfg.numFiles
	if fs.currentFileSize > fs.maxFileSize {
		fs.rotate()
	}
	// fs.logConfig()
}
