package c4solver

import (
	"os"
	"os/signal"

	log "github.com/igrek51/log15"
)

func HandleInterrupt(solver IMoveSolver) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Debug("Signal Interrupt - shutting down")
		solver.Interrupt()
	}()
}
