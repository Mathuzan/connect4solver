package main

import (
	"fmt"
	"sort"
	"time"

	log "github.com/inconshreveable/log15"
	"github.com/schollz/progressbar/v3"
)

type MoveSolver struct {
	cache              *EndingCache
	winner             Player
	lastBoardPrintTime time.Time
	progressBar        *progressbar.ProgressBar
	maxCacheDepth      uint

	iterations           uint64
	cacheUsages          uint64
	cachedDepthHistogram map[uint]uint64
	movesOrder           []int
}

const progressBarResolution = 1000000000

func NewMoveSolver(board *Board) *MoveSolver {
	maxCacheDepth := uint(24)

	log.Debug("Parameters set", log.Ctx{
		"maxCacheDepth": maxCacheDepth,
		"width":         board.w,
		"height":        board.h,
		"winStreak":     board.winStreak,
		"movesOrder":    CalculateMovesOrder(board),
	})

	return &MoveSolver{
		cache:                NewEndingCache(),
		lastBoardPrintTime:   time.Now(),
		progressBar:          progressbar.Default(progressBarResolution),
		cachedDepthHistogram: map[uint]uint64{},
		maxCacheDepth:        maxCacheDepth,
		movesOrder:           CalculateMovesOrder(board),
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

	keys := make([]uint, 0, len(s.cachedDepthHistogram))
	for k := range s.cachedDepthHistogram {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	for _, k := range keys {
		log.Debug(fmt.Sprintf("depth: %d, cached entries: %d", k, s.cachedDepthHistogram[k]))
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
) GameEnding {
	s.iterations++

	board.Throw(move, player)
	defer board.Revert(move)

	if depth <= s.maxCacheDepth {
		ending, ok := s.cache.Get(board)
		if ok {
			s.cacheUsages++
			return ending
		}
	}

	if s.iterations%10000 == 0 && time.Since(s.lastBoardPrintTime) >= 2*time.Second {
		s.lastBoardPrintTime = time.Now()
		s.ReportStatus(board, progressStart, progressEnd)
	}

	s.winner = board.HasWinner()
	if s.winner == PlayerA {
		return Win
	} else if s.winner == PlayerB {
		return Lose
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
					if depth <= s.maxCacheDepth {
						s.cache.Put(board, Win)
					}
					return Win
				}
				if bestEnding == nil || MoveResultsWeights[moveEnding] > MoveResultsWeights[*bestEnding] {
					bestEnding = &moveEnding
				}
			} else {
				// player B chooses worst possible move
				if moveEnding == Lose { // short-circuit, cant be worse
					if depth <= s.maxCacheDepth {
						s.cache.Put(board, Lose)
					}
					return Lose
				}
				if bestEnding == nil || MoveResultsWeights[moveEnding] < MoveResultsWeights[*bestEnding] {
					bestEnding = &moveEnding
				}
			}
		}
	}

	if bestEnding == nil {
		if depth <= s.maxCacheDepth {
			s.cache.Put(board, Tie)
		}
		return Tie
	}

	if depth <= s.maxCacheDepth {
		s.cache.Put(board, *bestEnding)
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

func (s *MoveSolver) ReportStatus(
	board *Board,
	progressStart float64,
	progressEnd float64,
) {
	log.Debug("Currently considered board", log.Ctx{
		"cacheSize":   s.cache.Size(),
		"iterations":  s.iterations,
		"cacheUsages": s.cacheUsages,
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
