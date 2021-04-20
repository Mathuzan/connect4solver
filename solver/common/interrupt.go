package common

import (
	"os"
	"os/signal"

	"github.com/pkg/errors"

	log "github.com/igrek51/log15"
)

type Interruptible interface {
	Interrupt()
}

func HandleInterrupt(interruptible Interruptible) chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Debug("Signal Interrupt - stopping")
		signal.Stop(c)
		interruptible.Interrupt()
	}()
	return c
}

var InterruptError error = errors.New("Interrupt")
