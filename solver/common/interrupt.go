package common

import (
	"os"
	"os/signal"

	"github.com/pkg/errors"

	log "github.com/igrek51/log15"
)

func HandleInterrupt(solver IMoveSolver) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Debug("Signal Interrupt - stopping")
		signal.Stop(c)
		solver.Interrupt()
	}()
}

var InterruptError error = errors.New("Interrupt")
