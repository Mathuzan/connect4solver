package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	log "github.com/igrek51/log15"
	"github.com/pkg/errors"
)

func main() {
	width, height, winStreak, profileEnabled, cacheEnabled := getArgs()

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
	solver := NewMoveSolver(board)

	if cacheEnabled && CacheFileExists(board) {
		cache, err := LoadCache(board)
		if err != nil {
			panic(errors.Wrap(err, "loading cache"))
		}
		solver.cache = cache
	}

	startTime := time.Now()
	endings := solver.MovesEndings(board)
	totalElapsed := time.Since(startTime)
	log.Info("Board solved", log.Ctx{
		"solveTime":   totalElapsed,
		"boardWidth":  width,
		"boardHeight": height,
		"winStreak":   winStreak,
	})
	for move, ending := range endings {
		playerEnding := EndingForPlayer(ending, myPlayer)
		log.Info(fmt.Sprintf("Best ending for move %d: %v", move, playerEnding))
	}

	if cacheEnabled {
		err := SaveCache(solver.cache)
		if err != nil {
			panic(errors.Wrap(err, "saving cache"))
		}
	}

	totalElapsed = time.Since(startTime)
	log.Info("Done", log.Ctx{
		"totalTime": totalElapsed,
	})
}
