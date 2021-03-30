package main

import (
	"fmt"
	"time"

	log "github.com/inconshreveable/log15"
	"github.com/schollz/progressbar/v3"
)

type MoveSolver struct {
	cache              *EndingCache
	winner             *Player
	lastBoardPrintTime time.Time
	progressBar        *progressbar.ProgressBar

	iterations       uint64
	cachedIterations uint64
}

const progressBarResolution = 1000000

func NewMoveSolver() *MoveSolver {
	return &MoveSolver{
		cache:              NewEndingCache(),
		lastBoardPrintTime: time.Now(),
		progressBar:        progressbar.Default(progressBarResolution),
	}
}

func (s *MoveSolver) BestEnding(board *Board) GameEnding {
	endings := s.MovesEndings(board)
	bestEnding := Lose
	for _, ending := range endings {
		if MoveResultsWeights[ending] > MoveResultsWeights[bestEnding] {
			bestEnding = ending
		}
	}
	return bestEnding
}

func (s *MoveSolver) MovesEndings(board *Board) []GameEnding {
	endings := make([]GameEnding, board.w)
	player := board.NextPlayer()

	for move := 0; move < board.w; move++ {
		progressStart := float64(move) / float64(board.w)
		progressEnd := float64(move+1) / float64(board.w)
		ending := s.bestEndingOnMove(board.Clone(), player, move, progressStart, progressEnd)
		endings[move] = ending
	}

	return endings
}

// bestEndingOnMove finds best ending on given next move
func (s *MoveSolver) bestEndingOnMove(
	board *Board,
	player Player,
	move int,
	progressStart float64,
	progressEnd float64,
) GameEnding {
	s.iterations++

	board.Throw(move, player)
	defer board.Revert(move)

	boardKey := s.cache.EvaluateKey(board)
	if s.cache.Has(boardKey) {
		s.cachedIterations++
		return s.cache.Get(boardKey)
	}

	if s.iterations%10000 == 0 && time.Since(s.lastBoardPrintTime) >= 2*time.Second {
		s.lastBoardPrintTime = time.Now()
		s.ReportStatus(board, progressStart, progressEnd)
	}

	s.winner = board.HasWinner()
	if s.winner != nil {
		if *s.winner == PlayerA {
			return Win
		} else {
			return Lose
		}
	}

	nextPlayer := oppositePlayer(player)

	// find further possible moves
	var bestEnding *GameEnding = nil
	for move := 0; move < board.w; move++ {
		if board.CanMakeMove(move) {
			moveEnding := s.bestEndingOnMove(board, nextPlayer, move,
				progressStart+float64(move)*(progressEnd-progressStart)/float64(board.w),
				progressStart+float64(move+1)*(progressEnd-progressStart)/float64(board.w),
			)

			if nextPlayer == PlayerA {
				// player A chooses highest possible move
				if moveEnding == Win { // short-circuit, cant be better
					s.cache.Put(boardKey, Win)
					return Win
				}
				if bestEnding == nil || MoveResultsWeights[moveEnding] > MoveResultsWeights[*bestEnding] {
					bestEnding = &moveEnding
				}
			} else {
				// player B chooses worst possible move
				if moveEnding == Lose { // short-circuit, cant be worse
					s.cache.Put(boardKey, Lose)
					return Lose
				}
				if bestEnding == nil || MoveResultsWeights[moveEnding] < MoveResultsWeights[*bestEnding] {
					bestEnding = &moveEnding
				}
			}
		}
	}

	if bestEnding == nil {
		s.cache.Put(boardKey, Tie)
		return Tie
	}

	s.cache.Put(boardKey, *bestEnding)
	return *bestEnding
}

func oppositePlayer(player Player) Player {
	if player == PlayerA {
		return PlayerB
	} else {
		return PlayerA
	}
}

func (s *MoveSolver) ReportStatus(
	board *Board,
	progressStart float64,
	progressEnd float64,
) {
	log.Debug("Currently considered board", log.Ctx{
		"cacheSize":        s.cache.Size(),
		"iterations":       s.iterations,
		"cachedIterations": s.cachedIterations,
	})
	fmt.Println(board.String())
	if s.progressBar != nil {
		s.progressBar.Set(int(progressStart * progressBarResolution))
	}
}

func (s *MoveSolver) BestEndingOnMove(
	board *Board,
	player Player,
	move int,
) GameEnding {
	return s.bestEndingOnMove(board, player, move, 0, 1)
}
