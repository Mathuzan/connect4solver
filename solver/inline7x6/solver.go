package inline7x6

import (
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
	W          int
	H          int
	tieDepth   uint

	startTime          time.Time
	lastBoardPrintTime time.Time
	firstProgress      float64
	lastProgress       float64
	progressBar        *progressbar.ProgressBar
	iterations         uint64
	lastIterations     uint64
	retrainMaxDepth    uint
}

func NewMoveSolver(board *common.Board) *MoveSolver {
	movesOrder := common.CalculateMovesOrder(board)
	cache := NewEndingCache(board.W, board.H)
	log.Debug("Solver configured", log.Ctx{
		"boardWidth":        board.W,
		"boardHeight":       board.H,
		"winStreak":         board.WinStreak,
		"movesOrder":        movesOrder,
		"maxCacheDepth":     cache.maxCachedDepth,
		"maxCacheDepthSize": cache.maxCacheDepthSize,
	})
	return &MoveSolver{
		W:                  board.W,
		H:                  board.H,
		cache:              cache,
		referee:            NewReferee(board),
		movesOrder:         movesOrder,
		lastBoardPrintTime: time.Now(),
		startTime:          time.Now(),
		progressBar:        common.NewProgressBar(),
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

	if depth <= 38 {
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
	if depth == 41 { // No more moves - Tie
		return common.Empty
	}

	// solve further possible moves of nextPlayer, at least one possible move is guaranteed
	nextPlayer := common.OppositePlayer(player)
	ties := 0
	for moveIndex := 0; moveIndex < 7; moveIndex++ {
		if board.CanMakeMove(s.movesOrder[moveIndex]) {
			moveEnding := s.bestEndingOnMove(board, nextPlayer, s.movesOrder[moveIndex],
				progressStart+float64(moveIndex)*(progressEnd-progressStart)/7.0,
				progressStart+float64(moveIndex+1)*(progressEnd-progressStart)/7.0,
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

func (s *MoveSolver) HasPlayerWon(board *common.Board, move int, y int, player common.Player) bool {
	return s.referee.HasPlayerWon(board, move, y, player)
}

func (s *MoveSolver) Interrupt() {
	s.interrupt = true
}

func (s *MoveSolver) Cache() common.ICache {
	return s.cache
}
