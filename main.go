package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	log "github.com/inconshreveable/log15"
)

func main() {
	width, height, profileEnabled := getArgs()

	if profileEnabled {
		log.Info("Starting CPU profiler")
		cpuProfile, _ := os.Create("cpuprof.prof")
		pprof.StartCPUProfile(cpuProfile)
		defer pprof.StopCPUProfile()
	}

	board := NewBoard(WithSize(width, height))
	fmt.Println(board.String())

	fmt.Println("Finding moves results...")
	startTime := time.Now()
	solver := NewMoveSolver(board)
	endings := solver.MovesEndings(board)
	for move, ending := range endings {
		log.Info(fmt.Sprintf("Best ending for move %d: %v", move, ending))
	}

	totalElapsed := time.Since(startTime)
	log.Info("Done", log.Ctx{"totalTime": totalElapsed})
}
