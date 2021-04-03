package main

import (
	"fmt"
	"time"

	log "github.com/igrek51/log15"
	"github.com/schollz/progressbar/v3"
)

type MoveSolver struct {
	cache              *EndingCache
	lastBoardPrintTime time.Time
	progressBar        *progressbar.ProgressBar

	iterations uint64
	movesOrder []int
}

const progressBarResolution = 1_000_000_000

func NewMoveSolver(board *Board) *MoveSolver {
	log.Debug("Parameters set", log.Ctx{
		"boardWidth":  board.w,
		"boardHeight": board.h,
		"winStreak":   board.winStreak,
		"movesOrder":  CalculateMovesOrder(board),
	})

	return &MoveSolver{
		cache:              NewEndingCache(board.w, board.h),
		lastBoardPrintTime: time.Now(),
		progressBar:        progressbar.Default(progressBarResolution),
		movesOrder:         CalculateMovesOrder(board),
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
		ending := s.bestEndingOnMove(board.Clone(), player, move, progressStart, progressEnd, 0)
		endings[move] = ending
	}

	//s.cache.ShowStatistics()

	return endings
}

// bestEndingOnMove finds best ending on given next move
func (s *MoveSolver) bestEndingOnMove(
	board *Board,
	player Player,
	move int,
	progressStart float64,
	progressEnd float64,
	depth uint,
) GameEnding {
	s.iterations++

	y := board.Throw(move, player)
	defer board.Revert(move, y)

	if depth <= s.cache.maxCacheDepth {
		ending, ok := s.cache.Get(board, depth)
		if ok {
			return ending
		}
	}

	if s.iterations%10000 == 0 && time.Since(s.lastBoardPrintTime) >= 2*time.Second {
		s.lastBoardPrintTime = time.Now()
		s.ReportStatus(board, depth, progressStart, progressEnd)
	}

	if board.referee.HasPlayerWon(board, move, y, player) {
		if player == PlayerA {
			return Win
		} else if player == PlayerB {
			return Lose
		}
	}

	nextPlayer := oppositePlayer(player)

	// find further possible moves
	var bestEnding *GameEnding = nil
	for moveIndex := 0; moveIndex < board.w; moveIndex++ {
		if board.CanMakeMove(s.movesOrder[moveIndex]) {
			moveEnding := s.bestEndingOnMove(board, nextPlayer, s.movesOrder[moveIndex],
				progressStart+float64(moveIndex)*(progressEnd-progressStart)/float64(board.w),
				progressStart+float64(moveIndex+1)*(progressEnd-progressStart)/float64(board.w),
				depth+1,
			)

			if nextPlayer == PlayerA {
				// player A chooses highest possible move
				if moveEnding == Win { // short-circuit, cant be better
					return s.cache.Put(board, depth, Win)
				}
				if bestEnding == nil || MoveResultsWeights[moveEnding] > MoveResultsWeights[*bestEnding] {
					bestEnding = &moveEnding
				}
			} else {
				// player B chooses worst possible move
				if moveEnding == Lose { // short-circuit, cant be worse
					return s.cache.Put(board, depth, Lose)
				}
				if bestEnding == nil || MoveResultsWeights[moveEnding] < MoveResultsWeights[*bestEnding] {
					bestEnding = &moveEnding
				}
			}
		}
	}

	if bestEnding == nil {
		return s.cache.Put(board, depth, Tie)
	}

	return s.cache.Put(board, depth, *bestEnding)
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
	depth uint,
	progressStart float64,
	progressEnd float64,
) {
	log.Debug("Currently considered board", log.Ctx{
		"cacheSize":   s.cache.Size(),
		"iterations":  s.iterations,
		"cacheUsages": s.cache.cacheUsages,
		"progress":    progressStart,
		"depth":       depth,
		"cacheClears": s.cache.clears,
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
	return s.bestEndingOnMove(board, player, move, 0, 1, 0)
}

func CalculateMovesOrder(board *Board) []int {
	movesOrder := []int{}
	pivot := board.w / 2
	for x := 0; x < board.w; x++ {
		var move int
		if x%2 == 0 {
			move = pivot - (x+1)/2
		} else {
			move = (pivot + (x+1)/2) % board.w
		}
		movesOrder = append(movesOrder, move)
	}
	return movesOrder
}
