package main

import (
	"time"

	"github.com/goccmack/goutil/log"
)

func log3() {
	for {
		log.Debug("Debug")
		log.Info("Info")
		time.Sleep(100 * time.Millisecond)
	}
}
