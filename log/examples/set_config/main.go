package main

import (
	"github.com/goccmack/goutil/log"
)

func main() {
	log.Info("Test started")

	// The following logging flows over two log files
	for i := 0; i < 2; i++ {
		log.Debugf("Debug %d", i)
		log.Infof("Debug %d", i)
	}

	// Change log file size and debug level for future messages
	log.SetConfig(3, 10000000, log.INFO)

	// Log more to the last log file.
	// Debug messages are suppressed
	for i := 2; i < 5; i++ {
		log.Debugf("Debug %d", i)
		log.Infof("Debug %d", i)
	}
	log.Info("Done")
}
