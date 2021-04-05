package c4solver

import (
	"fmt"
	"time"

	log "github.com/igrek51/log15"
	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"

	. "github.com/igrek51/connect4solver/c4solver/common"
)

type IMoveSolver interface {
	MovesEndings(board *Board) []Player
	HasPlayerWon(board *Board, move int, y int, player Player) bool
	Interrupt()
	PreloadCache(board *Board)
	SaveCache()
	ContextVars() log.Ctx
}

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

func NewMoveSolver(board *Board) *MoveSolver {
	movesOrder := CalculateMovesOrder(board)
	cache := NewEndingCache(board.W, board.H)
	log.Debug("Parameters set", log.Ctx{
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

func (s *MoveSolver) MovesEndings(board *Board) (endings []Player) {
	defer func() {
		if r := recover(); r != nil {
			_, ok := r.(InterruptType)
			if !ok {
				panic(r)
			}
			log.Debug("Interrupted")
			endings = nil
		}
	}()

	s.lastBoardPrintTime = time.Now()
	s.startTime = time.Now()
	s.iterations = 0
	endings = make([]Player, board.W)
	player := board.NextPlayer()

	for moveIndex := 0; moveIndex < board.W; moveIndex++ {
		move := s.movesOrder[moveIndex]
		progressStart := float64(moveIndex) / float64(board.W)
		progressEnd := float64(moveIndex+1) / float64(board.W)
		if board.CanMakeMove(move) {
			ending := s.bestEndingOnMove(board.Clone(), player, move, progressStart, progressEnd, 0)
			endings[move] = ending
		} else {
			endings[move] = NoMove
		}
	}

	return endings
}

const itReportPeriodMask = 0b11111111111111111111 // modulo 2^20 (1048576) mask

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

	if s.iterations&itReportPeriodMask == 0 && time.Since(s.lastBoardPrintTime) >= 2*time.Second {
		s.lastBoardPrintTime = time.Now()
		s.ReportStatus(board, progressStart, progressEnd)
		if s.interrupt {
			panic(InterruptError)
		}
	}

	if s.referee.HasPlayerWon(board, move, y, player) {
		return player
	}

	nextPlayer := oppositePlayer(player)

	// find further possible moves
	var bestEnding Player = Empty // Tie as a default ending (when no more moves)
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
	pivot := board.W / 2
	for x := 0; x < board.W; x++ {
		var move int
		if x%2 == 0 {
			move = pivot - (x+1)/2
		} else {
			move = (pivot + (x+1)/2) % board.W
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

func (s *MoveSolver) PreloadCache(board *Board) {
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

func (s *MoveSolver) ContextVars() log.Ctx {
	return log.Ctx{
		"cacheSize":   s.cache.Size(),
		"iterations":  s.iterations,
		"cacheUsages": s.cache.cacheUsages,
		"cacheClears": s.cache.clears,
	}
}

func (s *MoveSolver) HasPlayerWon(board *Board, move int, y int, player Player) bool {
	return s.referee.HasPlayerWon(board, move, y, player)
}