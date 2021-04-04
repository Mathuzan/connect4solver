package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	log "github.com/igrek51/log15"
)

var myPlayer = PlayerA

func main() {
	width, height, winStreak, profileEnabled, cacheEnabled := getArgs()

	if profileEnabled {
		log.Info("Starting CPU profiler")
		cpuProfile, _ := os.Create("cpuprof.prof")
		pprof.StartCPUProfile(cpuProfile)
		defer pprof.StopCPUProfile()
	}

	board := NewBoard(WithSize(width, height), WithWinStreak(winStreak))
	fmt.Println(board.String())

	fmt.Println("Finding moves results...")
	var solver IMoveSolver = NewMoveSolver(board)
	HandleInterrupt(solver)

	if cacheEnabled && CacheFileExists(board) {
		solver.PreloadCache(board)
	}

	startTime := time.Now()
	endings := solver.MovesEndings(board)
	totalElapsed := time.Since(startTime)

	logger := log.New(log.Ctx{
		"solveTime":   totalElapsed,
		"boardWidth":  width,
		"boardHeight": height,
		"winStreak":   winStreak,
	})
	logger.Info("Board solved", solver.ContextVars())
	for move, ending := range endings {
		playerEnding := EndingForPlayer(ending, myPlayer)
		log.Info(fmt.Sprintf("Best ending for move %d: %v", move, playerEnding))
	}

	if cacheEnabled {
		solver.SaveCache()
	}

	totalElapsed = time.Since(startTime)
	log.Info("Done", log.Ctx{
		"totalTime": totalElapsed,
	})
}
