package solver

import (
	"fmt"
	"time"

	log "github.com/igrek51/log15"
	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"

	"github.com/igrek51/connect4solver/solver/common"
)

type MoveSolver struct {
	cache              *EndingCache
	referee            *Referee
	lastBoardPrintTime time.Time
	startTime          time.Time
	progressBar        *progressbar.ProgressBar

	iterations uint64
	movesOrder []int
	interrupt  bool
}

const progressBarResolution = 1_000_000_000
const refreshProgressPeriod = 2000 * time.Millisecond
const itReportPeriodMask = 0b11111111111111111111 // modulo 2^20 (1048576) mask

func NewMoveSolver(board *common.Board) *MoveSolver {
	movesOrder := common.CalculateMovesOrder(board)
	cache := NewEndingCache(board.W, board.H)
	log.Debug("Solver configured", log.Ctx{
		"boardWidth":        board.W,
		"boardHeight":       board.H,
		"winStreak":         board.WinStreak,
		"movesOrder":        movesOrder,
		"maxCacheDepth":     cache.maxCacheDepth,
		"maxCacheDepthSize": cache.maxCacheDepthSize,
	})
	referee := NewReferee(board)

	return &MoveSolver{
		cache:              cache,
		lastBoardPrintTime: time.Now(),
		startTime:          time.Now(),
		progressBar:        progressbar.Default(progressBarResolution),
		movesOrder:         movesOrder,
		referee:            referee,
	}
}

func (s *MoveSolver) MovesEndings(board *common.Board) (endings []common.Player) {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok || !errors.Is(err, common.InterruptError) {
				panic(r)
			}
			log.Debug("Interrupted")
			endings = nil
		}
	}()

	s.lastBoardPrintTime = time.Now()
	s.startTime = time.Now()
	s.iterations = 0
	s.interrupt = false
	endings = make([]common.Player, board.W)
	player := board.NextPlayer()
	depth := board.CountMoves()

	for moveIndex := 0; moveIndex < board.W; moveIndex++ {
		move := s.movesOrder[moveIndex]
		progressStart := float64(moveIndex) / float64(board.W)
		progressEnd := float64(moveIndex+1) / float64(board.W)
		if board.CanMakeMove(move) {
			ending := s.bestEndingOnMove(board.Clone(), player, move, progressStart, progressEnd, depth)
			endings[move] = ending
		} else {
			endings[move] = common.NoMove
		}
	}

	return endings
}

// bestEndingOnMove finds best ending on given next move
func (s *MoveSolver) bestEndingOnMove(
	board *common.Board,
	player common.Player,
	move int,
	progressStart float64,
	progressEnd float64,
	depth uint,
) common.Player {
	s.iterations++

	y := board.Throw(move, player)
	defer board.Revert(move, y)

	if depth <= s.cache.maxCacheDepth {
		ending, ok := s.cache.Get(board, depth)
		if ok {
			s.cache.cacheUsages++
			return ending
		}
	}

	if s.iterations&itReportPeriodMask == 0 && time.Since(s.lastBoardPrintTime) >= refreshProgressPeriod {
		s.lastBoardPrintTime = time.Now()
		s.ReportStatus(board, progressStart, progressEnd)
		if s.interrupt {
			panic(common.InterruptError)
		}
	}

	if s.referee.HasPlayerWon(board, move, y, player) {
		return player
	}

	nextPlayer := common.OppositePlayer(player)

	// find further possible moves
	var bestEnding common.Player = common.Empty // Tie as a default ending (when no more moves)
	endingProcessed := 0
	for moveIndex := 0; moveIndex < board.W; moveIndex++ {
		if board.CanMakeMove(s.movesOrder[moveIndex]) {
			moveEnding := s.bestEndingOnMove(board, nextPlayer, s.movesOrder[moveIndex],
				progressStart+float64(moveIndex)*(progressEnd-progressStart)/float64(board.W),
				progressStart+float64(moveIndex+1)*(progressEnd-progressStart)/float64(board.W),
				depth+1,
			)

			if moveEnding == nextPlayer { // short-circuit, cant be better than winning
				return s.cache.Put(board, depth, moveEnding)
			}
			// player favors Tie over Lose
			if endingProcessed == 0 || moveEnding == common.Empty {
				bestEnding = moveEnding
			}
			endingProcessed++
		}
	}

	return s.cache.Put(board, depth, bestEnding)
}

// BestEndingOnMove finds best ending on given next move
func (s *MoveSolver) BestEndingOnMove(
	board *common.Board,
	player common.Player,
	move int,
) common.Player {
	return s.bestEndingOnMove(board, player, move, 0, 1, 0)
}

func (s *MoveSolver) HasPlayerWon(board *common.Board, move int, y int, player common.Player) bool {
	return s.referee.HasPlayerWon(board, move, y, player)
}

func (s *MoveSolver) PreloadCache(board *common.Board) {
	cache, err := LoadCache(board)
	if err != nil {
		panic(errors.Wrap(err, "loading cache"))
	}
	s.cache = cache
}

func (s *MoveSolver) SaveCache() {
	err := SaveCache(s.cache)
	if err != nil {
		panic(errors.Wrap(err, "saving cache"))
	}
}

func (s *MoveSolver) Interrupt() {
	s.interrupt = true
}

func (s *MoveSolver) ReportStatus(
	board *common.Board,
	progressStart float64,
	progressEnd float64,
) {
	duration := time.Since(s.startTime)
	var eta time.Duration
	if progressStart > 0 && duration > 0 {
		eta = time.Duration((1 - progressStart) / (progressStart / float64(duration)))
	}
	iterationsPerSec := common.BigintSeparated(s.iterations * uint64(time.Second) / uint64(duration))

	log.Debug("Currently considered board", log.Ctx{
		"cacheSize":         common.BigintSeparated(s.cache.Size()),
		"iterations":        common.BigintSeparated(s.iterations),
		"cacheClears":       common.BigintSeparated(s.cache.clears),
		"maxUnclearedDepth": maximumZeroIndex(s.cache.depthClears),
		"progress":          fmt.Sprintf("%v", progressStart),
		"eta":               eta,
		"ips":               iterationsPerSec,
	})
	fmt.Println(board.String())
	if s.progressBar != nil {
		s.progressBar.Set(int(progressStart * progressBarResolution))
	}
}

func (s *MoveSolver) SummaryVars() log.Ctx {
	return log.Ctx{
		"cacheSize":   common.BigintSeparated(s.cache.Size()),
		"iterations":  common.BigintSeparated(s.iterations),
		"cacheClears": common.BigintSeparated(s.cache.clears),
	}
}

func maximumZeroIndex(nums []uint64) int {
	maxi := -1
	for i, num := range nums {
		if num == 0 {
			maxi = i
		} else {
			break
		}
	}
	return maxi
}
