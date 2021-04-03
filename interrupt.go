package main

import (
	"os"
	"os/signal"

	"github.com/pkg/errors"

	log "github.com/igrek51/log15"
)

func HandleInterrupt(solver *MoveSolver) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Debug("Signal Interrupt")
		solver.Interrupt = true
	}()
}

type InterruptType error

var InterruptError InterruptType = errors.New("Interrupt")
