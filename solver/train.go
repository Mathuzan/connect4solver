package solver

import (
	"fmt"
	"time"

	log "github.com/igrek51/log15"

	"github.com/igrek51/connect4solver/solver/common"
)

func Train(width, height, winStreak int, cacheEnabled bool) {
	board := common.NewBoard(common.WithSize(width, height), common.WithWinStreak(winStreak))
	fmt.Println(board.String())

	solver := CreateSolver(board)
	if cacheEnabled && CacheFileExists(board) {
		solver.PreloadCache(board)
	}
	common.HandleInterrupt(solver)

	startTime := time.Now()
	endings := solver.MovesEndings(board)
	totalElapsed := time.Since(startTime)

	logger := log.New(log.Ctx{
		"solveTime":   totalElapsed,
		"boardWidth":  width,
		"boardHeight": height,
		"winStreak":   winStreak,
	})
	logger.Info("Board solved", solver.SummaryVars())

	player := board.NextPlayer()
	for move, ending := range endings {
		if ending != common.NoMove {
			playerEnding := common.EndingForPlayer(ending, player)
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
