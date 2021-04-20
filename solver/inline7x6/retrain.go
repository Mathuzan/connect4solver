package inline7x6

import (
	"os/signal"
	"time"

	log "github.com/igrek51/log15"
	"github.com/pkg/errors"

	"github.com/igrek51/connect4solver/solver/common"
)

func (s *MoveSolver) Retrain(board *common.Board, maxDepth uint) {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok || !errors.Is(err, common.InterruptError) {
				panic(r)
			}
			log.Debug("Interrupted")
		}
	}()
	interruptChannel := common.HandleInterrupt(s)

	s.startTime = time.Now()
	s.lastBoardPrintTime = time.Now()
	s.firstProgress = 0
	s.iterations = 0
	s.interrupt = false

	depth := board.CountMoves()
	player := board.NextPlayer()

	for moveIndex := 0; moveIndex < board.W; moveIndex++ {
		move := s.movesOrder[moveIndex]
		progressStart := float64(moveIndex) / float64(board.W)
		progressEnd := float64(moveIndex+1) / float64(board.W)
		if board.CanMakeMove(move) {
			s.retrainEndingOnMove(board, player, move, progressStart, progressEnd, depth, maxDepth)
		}
	}

	signal.Stop(interruptChannel)
}

// solve board without short-circuit features
func (s *MoveSolver) retrainEndingOnMove(
	board *common.Board,
	player common.Player,
	move int,
	progressStart float64,
	progressEnd float64,
	depth uint,
	maxDepth uint,
) common.Player {
	s.iterations++

	y := board.Throw(move, player)
	defer board.Revert(move, y)

	if depth >= maxDepth && depth <= 38 {
		ending, ok := s.cache.Get(board, depth)
		if ok {
			s.cache.cacheUsages++
			return ending
		}
	}

	s.reportCycle(board, progressStart)

	if s.referee.HasPlayerWon(board, move, y, player) {
		return player
	}
	if depth == s.tieDepth { // No more moves - Tie
		return common.Empty
	}

	// solve further possible moves of nextPlayer, at least one possible move is guaranteed
	nextPlayer := common.OppositePlayer(player)
	wins := 0
	ties := 0
	var moveEnding common.Player
	for moveIndex := 0; moveIndex < 7; moveIndex++ {
		if board.CanMakeMove(s.movesOrder[moveIndex]) {
			if depth+1 < maxDepth {
				moveEnding = s.retrainEndingOnMove(board, nextPlayer, s.movesOrder[moveIndex],
					progressStart+float64(moveIndex)*(progressEnd-progressStart)/7.0,
					progressStart+float64(moveIndex+1)*(progressEnd-progressStart)/7.0,
					depth+1, maxDepth,
				)
			} else {
				moveEnding = s.bestEndingOnMove(board, nextPlayer, s.movesOrder[moveIndex],
					progressStart+float64(moveIndex)*(progressEnd-progressStart)/7.0,
					progressStart+float64(moveIndex+1)*(progressEnd-progressStart)/7.0,
					depth+1,
				)
			}

			if moveEnding == nextPlayer {
				wins++
			} else if moveEnding == common.Empty {
				ties++
			}
		}
	}
	if wins > 0 {
		return s.cache.Put(board, depth, nextPlayer)
	} else if ties > 0 {
		return s.cache.Put(board, depth, common.Empty)
	} else {
		return s.cache.Put(board, depth, player)
	}
}
