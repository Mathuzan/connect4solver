package main

import (
	"fmt"
	"time"

	log "github.com/igrek51/log15"
)

func Train(width, height, winStreak int, cacheEnabled bool) {
	board := NewBoard(WithSize(width, height), WithWinStreak(winStreak))
	fmt.Println(board.String())

	log.Debug("Finding moves results...")
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

	player := board.NextPlayer()
	for move, ending := range endings {
		if ending != NoMove {
			playerEnding := EndingForPlayer(ending, player)
			log.Info(fmt.Sprintf("Best ending for move %d: %v", move, playerEnding))
		}
	}

	if cacheEnabled {
		solver.SaveCache()
	}

	totalElapsed = time.Since(startTime)
	log.Info("Done", log.Ctx{
		"totalTime": totalElapsed,
	})
}
