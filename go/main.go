package main

import (
	"fmt"
	"time"

	log "github.com/inconshreveable/log15"
)

func main() {
	width, height := getArgs()
	board := NewBoard(WithSize(width, height))
	fmt.Println(board.String())

	fmt.Println("Finding moves results...")
	startTime := time.Now()
	solver := NewMoveSolver()
	endings := solver.MovesEndings(board)
	for move, ending := range endings {
		log.Info("Move ending", log.Ctx{"move": move, "result": ending})
	}

	totalElapsed := time.Since(startTime)
	log.Info("Done", log.Ctx{"totalTime": totalElapsed})
}
