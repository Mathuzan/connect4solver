package inline5x5

import (
	"math/bits"

	"github.com/igrek51/connect4solver/c4solver/common"
)

type Referee struct {
	w         int
	h         int
	winStreak int

	verticalMovesMap []common.Player
	binaryRowMap     []bool
	winStreak1       int
}

func NewReferee(board *common.Board) *Referee {
	s := &Referee{
		w:          board.W,
		h:          board.H,
		winStreak:  board.WinStreak,
		winStreak1: board.WinStreak - 1,
	}

	// get all possible column layouts and pre-calculate winners for them
	verticalMovesMap := make([]common.Player, 1<<(s.h+1))
	for colState := uint64(0); colState <= (1<<(s.h+1))-1; colState++ {
		verticalMovesMap[colState] = s.whoWonColumn(colState)
	}
	s.verticalMovesMap = verticalMovesMap

	maxSize := s.w
	if s.h > maxSize {
		maxSize = s.h
	}
	// binary row is represented as: 1 - desired token of player, 0 - opponent's token or empty
	binaryRowMap := make([]bool, 1<<maxSize)
	for binaryRow := uint64(0); binaryRow <= (1<<maxSize)-1; binaryRow++ {
		binaryRowMap[binaryRow] = s.hasWonRow(binaryRow)
	}
	s.binaryRowMap = binaryRowMap

	return s
}

func (s *Referee) HasPlayerWon(board *common.Board, move int, y int, player common.Player) bool {
	return s.HasPlayerWonVertical(board, move, player) ||
		s.HasPlayerWonHorizontal(board, y, player) ||
		s.HasPlayerWonDiagonal(board, move, y, player) ||
		s.HasPlayerWonDiagonalCounter(board, move, y, player)
}

func (s *Referee) HasPlayerWonVertical(board *common.Board, move int, player common.Player) bool {
	return s.verticalMovesMap[board.State[move]] == player
}

func (s *Referee) HasPlayerWonHorizontal(board *common.Board, y int, player common.Player) bool {
	var binaryRow uint64
	for x := 0; x < 5; x++ {
		if board.GetCell(x, y) == player {
			binaryRow |= 1 << x
		}
	}
	return s.binaryRowMap[binaryRow]
}

func (s *Referee) HasPlayerWonDiagonal(board *common.Board, startX int, startY int, player common.Player) bool {
	var binaryRow uint64
	y := startY - 3
	for x := startX - 3; x <= startX+3; x++ {
		if x >= 0 && x < 5 && y >= 0 && y < 5 {
			if board.GetCell(x, y) == player {
				binaryRow |= 1 << x
			}
		}
		y++
	}
	return s.binaryRowMap[binaryRow]
}

func (s *Referee) HasPlayerWonDiagonalCounter(board *common.Board, startX int, startY int, player common.Player) bool {
	var binaryRow uint64
	y := startY + 3
	for x := startX - 3; x <= startX+3; x++ {
		if x >= 0 && x < 5 && y >= 0 && y < 5 {
			if board.GetCell(x, y) == player {
				binaryRow |= 1 << x
			}
		}
		y--
	}
	return s.binaryRowMap[binaryRow]
}

func getStackSize(columnState uint64) int {
	return 7 - bits.LeadingZeros8(uint8(columnState))
}

func (s *Referee) whoWonColumn(columnState uint64) common.Player {
	stackSize := getStackSize(columnState)
	if stackSize < s.winStreak {
		return common.Empty
	}

	onesB := columnState
	onesA := ^columnState

	for i := 0; i < s.winStreak1; i++ {
		onesB &= onesB >> 1
		onesA &= onesA >> 1
	}

	var mask uint64 = (1 << (stackSize - s.winStreak1)) - 1

	if onesA&mask != 0 {
		return common.PlayerA
	}
	if onesB&mask != 0 {
		return common.PlayerB
	}
	return common.Empty
}

func (s *Referee) hasWonRow(row uint64) bool {
	ones := row
	for i := 0; i < s.winStreak1; i++ {
		ones &= ones >> 1
	}
	return ones != 0
}
