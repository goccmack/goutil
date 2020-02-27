package main

import (
	"time"

	"github.com/goccmack/goutil/log/examples/basic/pkga"

	"github.com/goccmack/goutil/log"
)

func main() {
	// Log something
	log.Info("This message WILL appear in the log")

	pkga.Go()
	// Give pkga time to log
	time.Sleep(time.Millisecond)

	log.Debug("This message will NOT appear in the log")
	log.Panic("This is a panic")
}
