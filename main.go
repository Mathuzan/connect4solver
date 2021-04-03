package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	log "github.com/igrek51/log15"
)

func main() {
	width, height, winStreak, profileEnabled := getArgs()

	if profileEnabled {
		log.Info("Starting CPU profiler")
		cpuProfile, _ := os.Create("cpuprof.prof")
		pprof.StartCPUProfile(cpuProfile)
		defer pprof.StopCPUProfile()
	}

	board := NewBoard(WithSize(width, height), WithWinStreak(winStreak))
	myPlayer := PlayerA
	fmt.Println(board.String())

	fmt.Println("Finding moves results...")
	startTime := time.Now()
	solver := NewMoveSolver(board)
	endings := solver.MovesEndings(board)
	for move, ending := range endings {
		playerEnding := EndingForPlayer(ending, myPlayer)
		log.Info(fmt.Sprintf("Best ending for move %d: %v", move, playerEnding))
	}

	totalElapsed := time.Since(startTime)
	log.Info("Done", log.Ctx{
		"totalTime":   totalElapsed,
		"boardWidth":  width,
		"boardHeight": height,
		"winStreak":   winStreak,
	})
}
