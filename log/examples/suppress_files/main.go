package main

import (
	"time"

	"github.com/goccmack/goutil/log"
)

func main() {
	// pkga and pkgb log messages are interleaved
	go log1()
	go log2()
	go log3()

	// Let them log for a second
	time.Sleep(time.Second)

	// Suppress pkga and pkgb debug messages. Their info message are still present.
	// ".go" file extension is optional.
	log.Suppress("file1.go,file2")

	// Let them log for a second
	time.Sleep(time.Second)

	log.Info("Done")
}
