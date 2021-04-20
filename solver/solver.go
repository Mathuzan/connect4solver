package solver

import (
	"fmt"
	"os/signal"
	"time"

	log "github.com/igrek51/log15"
	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"

	"github.com/igrek51/connect4solver/solver/common"
)

type MoveSolver struct {
	cache      *EndingCache
	referee    *Referee
	movesOrder []int
	interrupt  bool

	tieDepth uint

	startTime          time.Time
	lastBoardPrintTime time.Time
	firstProgress      float64
	progressBar        *progressbar.ProgressBar
	iterations         uint64
	lastIterations     uint64
}

const progressBarResolution = 1_000_000_000
const refreshProgressPeriod = 2 * time.Second
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
		interrupt:          false,
		tieDepth:           uint(board.W*board.H - 1),
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
	interruptChannel := common.HandleInterrupt(s)

	s.startTime = time.Now()
	s.lastBoardPrintTime = time.Now()
	s.firstProgress = 0
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

	signal.Stop(interruptChannel)
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

	s.reportCycle(board, progressStart)

	if s.referee.HasPlayerWon(board, move, y, player) {
		return player
	}
	if depth == s.tieDepth { // No more moves - Tie
		return common.Empty
	}

	// solve further possible moves of nextPlayer, at least one possible move is guaranteed
	nextPlayer := common.OppositePlayer(player)
	ties := 0
	for moveIndex := 0; moveIndex < board.W; moveIndex++ {
		if board.CanMakeMove(s.movesOrder[moveIndex]) {
			moveEnding := s.bestEndingOnMove(board, nextPlayer, s.movesOrder[moveIndex],
				progressStart+float64(moveIndex)*(progressEnd-progressStart)/float64(board.W),
				progressStart+float64(moveIndex+1)*(progressEnd-progressStart)/float64(board.W),
				depth+1,
			)

			if moveEnding == nextPlayer { // short-circuit, cant be better than winning
				return s.cache.Put(board, depth, nextPlayer)
			}
			if moveEnding == common.Empty {
				ties++
			}
		}
	}
	// player favors Tie over Lose
	if ties > 0 {
		return s.cache.Put(board, depth, common.Empty)
	} else {
		return s.cache.Put(board, depth, player)
	}
}

// BestEndingOnMove finds best ending on given next move
func (s *MoveSolver) BestEndingOnMove(
	board *common.Board,
	player common.Player,
	move int,
) common.Player {
	depth := board.CountMoves()
	return s.bestEndingOnMove(board, player, move, 0, 1, depth)
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

func (s *MoveSolver) reportCycle(board *common.Board, progressStart float64) {
	if s.iterations&itReportPeriodMask == 0 && time.Since(s.lastBoardPrintTime) >= refreshProgressPeriod {
		s.ReportStatus(board, progressStart)
		if s.interrupt {
			panic(common.InterruptError)
		}
	}
}

func (s *MoveSolver) ReportStatus(
	board *common.Board,
	progress float64,
) {
	duration := time.Since(s.startTime)
	instDuration := time.Since(s.lastBoardPrintTime)
	if s.firstProgress == 0 {
		s.firstProgress = progress
	}
	var eta time.Duration
	if progress > s.firstProgress && duration > 0 {
		eta = time.Duration((1-progress)*1000/(progress-s.firstProgress)) * (duration - refreshProgressPeriod) / 1000
	}
	iterationsPerSec := common.BigintSeparated(s.iterations / uint64(duration/time.Second))
	instIterationsPerSec := common.BigintSeparated((s.iterations - s.lastIterations) / uint64(instDuration/time.Second))
	s.lastIterations = s.iterations
	s.lastBoardPrintTime = time.Now()

	firstCaches := []int{}
	for d := uint(0); d < uint(10); d++ {
		firstCaches = append(firstCaches, s.Cache().DepthSize(d))
	}

	log.Debug("Currently considered board", log.Ctx{
		"cacheSize":         common.BigintSeparated(s.cache.Size()),
		"iterations":        common.BigintSeparated(s.iterations),
		"cacheClears":       common.BigintSeparated(s.cache.clears),
		"maxUnclearedDepth": maximumZeroIndex(s.cache.depthClears),
		"progress":          fmt.Sprintf("%v", progress),
		"eta":               eta.Truncate(time.Second),
		"ipsAvg":            iterationsPerSec,
		"ips":               instIterationsPerSec,
		"firstCacheLen":     firstCaches,
	})
	fmt.Println(board.String())
	if s.progressBar != nil {
		s.progressBar.Set(int(progress * progressBarResolution))
	}
}

func (s *MoveSolver) SummaryVars() log.Ctx {
	return log.Ctx{
		"cacheSize":   common.BigintSeparated(s.cache.Size()),
		"iterations":  common.BigintSeparated(s.iterations),
		"cacheClears": common.BigintSeparated(s.cache.clears),
	}
}

func (s *MoveSolver) Cache() common.ICache {
	return s.cache
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

	if depth >= maxDepth && depth <= s.cache.maxCacheDepth {
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
	for moveIndex := 0; moveIndex < board.W; moveIndex++ {
		if board.CanMakeMove(s.movesOrder[moveIndex]) {
			if depth+1 < maxDepth {
				moveEnding = s.retrainEndingOnMove(board, nextPlayer, s.movesOrder[moveIndex],
					progressStart+float64(moveIndex)*(progressEnd-progressStart)/float64(board.W),
					progressStart+float64(moveIndex+1)*(progressEnd-progressStart)/float64(board.W),
					depth+1, maxDepth,
				)
			} else {
				moveEnding = s.bestEndingOnMove(board, nextPlayer, s.movesOrder[moveIndex],
					progressStart+float64(moveIndex)*(progressEnd-progressStart)/float64(board.W),
					progressStart+float64(moveIndex+1)*(progressEnd-progressStart)/float64(board.W),
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
