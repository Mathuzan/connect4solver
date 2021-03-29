package main

import (
	"fmt"
	"time"

	log "github.com/inconshreveable/log15"
)

type MoveSolver struct {
	cache *EndingCache
}

func NewMoveSolver() *MoveSolver {
	return &MoveSolver{
		cache: NewEndingCache(),
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
		ending := s.BestEndingOnMove(board.Clone(), player, move)
		endings[move] = ending
	}

	return endings
}

var winner *Player
var lastBoardPrintTime time.Time

// BestEndingOnMove finds best ending on given next move
func (s *MoveSolver) BestEndingOnMove(board *Board, player Player, move int) GameEnding {
	board.Throw(move, player)
	defer board.Revert(move)

	boardKey := s.cache.EvaluateKey(board)
	if s.cache.Has(boardKey) {
		return s.cache.Get(boardKey)
	}

	if time.Since(lastBoardPrintTime) >= 2*time.Second {
		lastBoardPrintTime = time.Now()
		ReportStatus(board, s.cache)
	}

	winner = board.HasWinner()
	if winner != nil {
		if *winner == PlayerA {
			return Win
		} else {
			return Lose
		}
	}

	ending := s.NextMoveEnding(board, oppositePlayer(player))
	s.cache.Put(boardKey, ending)
	return ending
}

// NextMoveEnding finds further possible moves
func (s *MoveSolver) NextMoveEnding(board *Board, player Player) GameEnding {
	var bestEnding *GameEnding = nil
	for move := 0; move < board.w; move++ {
		if board.CanMakeMove(move) {
			moveEnding := s.BestEndingOnMove(board, player, move)

			if player == PlayerA {
				// player A chooses highest possible move
				if moveEnding == Win { // short-circuit, cant be better
					return Win
				}
				if bestEnding == nil || MoveResultsWeights[moveEnding] > MoveResultsWeights[*bestEnding] {
					bestEnding = &moveEnding
				}
			} else {
				// player B chooses worst possible move
				if moveEnding == Lose { // short-circuit, cant be worse
					return Lose
				}
				if bestEnding == nil || MoveResultsWeights[moveEnding] < MoveResultsWeights[*bestEnding] {
					bestEnding = &moveEnding
				}
			}
		}
	}

	if bestEnding == nil {
		return Tie
	}

	return *bestEnding
}

func oppositePlayer(player Player) Player {
	if player == PlayerA {
		return PlayerB
	} else {
		return PlayerA
	}
}

func maxPossibleMove(endings []GameEnding) GameEnding {
	maxr := endings[0]
	for _, ending := range endings {
		if MoveResultsWeights[ending] > MoveResultsWeights[maxr] {
			maxr = ending
		}
	}
	return maxr
}

func minPossibleMove(endings []GameEnding) GameEnding {
	minr := endings[0]
	for _, ending := range endings {
		if MoveResultsWeights[ending] < MoveResultsWeights[minr] {
			minr = ending
		}
	}
	return minr
}

func ReportStatus(board *Board, cache *EndingCache) {
	log.Debug("Currently considered board", log.Ctx{
		"cacheSize": cache.Size(),
	})
	fmt.Println(board.String())
}
