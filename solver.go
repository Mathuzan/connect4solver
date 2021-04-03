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
	startTime          time.Time
	progressBar        *progressbar.ProgressBar

	iterations uint64
	movesOrder []int
}

const progressBarResolution = 1_000_000_000

func NewMoveSolver(board *Board) *MoveSolver {
	movesOrder := CalculateMovesOrder(board)
	cache := NewEndingCache(board.w, board.h)
	log.Debug("Parameters set", log.Ctx{
		"boardWidth":        board.w,
		"boardHeight":       board.h,
		"winStreak":         board.winStreak,
		"movesOrder":        movesOrder,
		"maxCacheDepth":     cache.maxCacheDepth,
		"maxCacheDepthSize": cache.maxCacheDepthSize,
	})

	return &MoveSolver{
		cache:              cache,
		lastBoardPrintTime: time.Now(),
		startTime:          time.Now(),
		progressBar:        progressbar.Default(progressBarResolution),
		movesOrder:         movesOrder,
	}
}

func (s *MoveSolver) MovesEndings(board *Board) []Player {
	endings := make([]Player, board.w)
	player := board.NextPlayer()

	for move := 0; move < board.w; move++ {
		progressStart := float64(move) / float64(board.w)
		progressEnd := float64(move+1) / float64(board.w)
		ending := s.bestEndingOnMove(board.Clone(), player, move, progressStart, progressEnd, 0)
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
	depth uint,
) Player {
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
		s.ReportStatus(board, progressStart, progressEnd)
	}

	if board.referee.HasPlayerWon(board, move, y, player) {
		return player
	}

	nextPlayer := oppositePlayer(player)

	// find further possible moves
	var bestEnding Player = Empty // Tie as a default ending (when no more moves)
	endingProcessed := 0
	for moveIndex := 0; moveIndex < board.w; moveIndex++ {
		if board.CanMakeMove(s.movesOrder[moveIndex]) {
			moveEnding := s.bestEndingOnMove(board, nextPlayer, s.movesOrder[moveIndex],
				progressStart+float64(moveIndex)*(progressEnd-progressStart)/float64(board.w),
				progressStart+float64(moveIndex+1)*(progressEnd-progressStart)/float64(board.w),
				depth+1,
			)

			if moveEnding == nextPlayer { // short-circuit, cant be better than winning
				return s.cache.Put(board, depth, moveEnding)
			}
			// player favors Tie over Lose
			if endingProcessed == 0 || moveEnding == Empty {
				bestEnding = moveEnding
			}
			endingProcessed++
		}
	}

	return s.cache.Put(board, depth, bestEnding)
}

func oppositePlayer(player Player) Player {
	return 1 - player
}

func (s *MoveSolver) ReportStatus(
	board *Board,
	progressStart float64,
	progressEnd float64,
) {
	duration := time.Since(s.startTime)
	var eta time.Duration
	if progressStart > 0 && duration > 0 {
		eta = time.Duration((1 - progressStart) / (progressStart / float64(duration)))
	}

	log.Debug("Currently considered board", log.Ctx{
		"cacheSize":   s.cache.Size(),
		"iterations":  s.iterations,
		"cacheUsages": s.cache.cacheUsages,
		"progress":    progressStart,
		"cacheClears": s.cache.clears,
		"eta":         eta,
	})
	fmt.Println(board.String())
	if s.progressBar != nil {
		s.progressBar.Set(int(progressStart * progressBarResolution))
	}
}

// BestEndingOnMove finds best ending on given next move
func (s *MoveSolver) BestEndingOnMove(
	board *Board,
	player Player,
	move int,
) Player {
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

func EndingForPlayer(ending Player, player Player) GameEnding {
	if ending == Empty {
		return Tie
	}
	if ending == player {
		return Win
	} else {
		return Lose
	}
}
