package main

import (
	"github.com/goccmack/goutil/log"
)

func main() {
	log.Exitf(2, "Logging exit code %d", 2)
}
