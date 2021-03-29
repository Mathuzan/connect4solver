package main

import (
	"fmt"
	"time"

	log "github.com/inconshreveable/log15"
)

type MoveSolver struct {
	cache *EndingStore
}

func NewMoveSolver() *MoveSolver {
	return &MoveSolver{
		cache: NewEndingStore(),
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
	possibleMovesEndings := make([]GameEnding, 0, board.w)
	for move := 0; move < board.w; move++ {
		if board.CanMakeMove(move) {
			moveEnding := s.BestEndingOnMove(board, player, move)
			if player == PlayerA && moveEnding == Win {
				return Win
			}
			if player == PlayerB && moveEnding == Lose {
				return Lose
			}

			possibleMovesEndings = append(possibleMovesEndings, moveEnding)
		}
	}

	if len(possibleMovesEndings) == 0 {
		return Tie
	}

	if player == PlayerA {
		return maxPossibleMove(possibleMovesEndings)
	} else {
		return minPossibleMove(possibleMovesEndings)
	}
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

func ReportStatus(board *Board, cache *EndingStore) {
	log.Debug("Currently considered board", log.Ctx{
		"cacheSize": cache.Size(),
	})
	fmt.Println(board.String())
}
