package c4solver

import (
	"fmt"
	"time"

	log "github.com/igrek51/log15"

	"github.com/igrek51/connect4solver/c4solver/common"
	"github.com/igrek51/connect4solver/c4solver/inline5x5"
	"github.com/igrek51/connect4solver/c4solver/inline7x6"
)

func Train(width, height, winStreak int, cacheEnabled bool) {
	board := common.NewBoard(common.WithSize(width, height), common.WithWinStreak(winStreak))
	fmt.Println(board.String())

	log.Debug("Finding moves results...")
	solver := createSolver(board)

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
	logger.Info("Board solved", solver.SummaryVars())

	player := board.NextPlayer()
	for move, ending := range endings {
		if ending != common.NoMove {
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

func createSolver(board *common.Board) IMoveSolver {
	var solver IMoveSolver
	// take precedence with inlined optimized solvers
	if board.W == 7 && board.H == 6 {
		solver = inline7x6.NewMoveSolver(board)
	} else if board.W == 5 && board.H == 5 {
		solver = inline5x5.NewMoveSolver(board)
	} else {
		solver = NewMoveSolver(board)
	}
	HandleInterrupt(solver)
	return solver
}
