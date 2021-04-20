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
	if cacheEnabled && common.CacheFileExists(board) {
		common.MustLoadCache(solver.Cache(), board.W, board.H)
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
	logger.Info("Board solved", solver.SummaryVars())

	player := board.NextPlayer()
	for move, ending := range endings {
		if ending != common.NoMove {
			playerEnding := common.EndingForPlayer(ending, player)
			log.Info(fmt.Sprintf("Best ending for move %d: %v", move, playerEnding))
		}
	}

	fmt.Println(board.String())
	printEndingsLine(endings, player)

	if cacheEnabled {
		common.MustSaveCache(solver.Cache(), board.W, board.H)
	}

	totalElapsed = time.Since(startTime)
	log.Info("Done", log.Ctx{
		"totalTime": totalElapsed,
	})
}
