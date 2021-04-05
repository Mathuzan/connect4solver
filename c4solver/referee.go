package c4solver

import (
	"math/bits"

	. "github.com/igrek51/connect4solver/c4solver/common"
)

type Referee struct {
	w         int
	h         int
	winStreak int

	verticalMovesMap []Player
	binaryRowMap     []bool
	winStreak1       int

	winner       Player
	lastToken    Player
	currentToken Player
	sameStreak   int
	stacksSum    uint64
}

func NewReferee(board *Board) *Referee {
	s := &Referee{
		w:          board.W,
		h:          board.H,
		winStreak:  board.WinStreak,
		winStreak1: board.WinStreak - 1,
	}

	// get all possible column layouts and pre-calculate winners for them
	verticalMovesMap := make([]Player, 1<<(s.h+1))
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

func (s *Referee) HasPlayerWon(board *Board, move int, y int, player Player) bool {
	return s.HasPlayerWonVertical(board, move, player) ||
		s.HasPlayerWonHorizontal(board, y, player) ||
		s.HasPlayerWonDiagonal(board, move, y, player) ||
		s.HasPlayerWonDiagonalCounter(board, move, y, player)
}

func (s *Referee) HasPlayerWonVertical(board *Board, move int, player Player) bool {
	return s.verticalMovesMap[board.State[move]] == player
}

func (s *Referee) HasPlayerWonHorizontal(board *Board, y int, player Player) bool {
	var binaryRow uint64
	for x := 0; x < s.w; x++ {
		if board.GetCell(x, y) == player {
			binaryRow |= 1 << x
		}
	}
	return s.binaryRowMap[binaryRow]
}

func (s *Referee) HasPlayerWonDiagonal(board *Board, startX int, startY int, player Player) bool {
	var binaryRow uint64
	y := startY - s.winStreak1
	for x := startX - s.winStreak1; x <= startX+s.winStreak1; x++ {
		if x >= 0 && x < s.w && y >= 0 && y <= s.h {
			if board.GetCell(x, y) == player {
				binaryRow |= 1 << x
			}
		}
		y++
	}
	return s.binaryRowMap[binaryRow]
}

func (s *Referee) HasPlayerWonDiagonalCounter(board *Board, startX int, startY int, player Player) bool {
	var binaryRow uint64
	y := startY + s.winStreak1
	for x := startX - s.winStreak1; x <= startX+s.winStreak1; x++ {
		if x >= 0 && x < s.w && y >= 0 && y <= s.h {
			if board.GetCell(x, y) == player {
				binaryRow |= 1 << x
			}
		}
		y--
	}
	return s.binaryRowMap[binaryRow]
}

func (s *Referee) HasWinner(board *Board) Player {
	s.winner = s.checkVertical(board)
	if s.winner != Empty {
		return s.winner
	}

	s.winner = s.checkHorizontal(board)
	if s.winner != Empty {
		return s.winner
	}

	s.winner = s.checkDiagonals(board)
	if s.winner != Empty {
		return s.winner
	}

	return Empty
}

func (s *Referee) checkVertical(board *Board) Player {
	for x := 0; x < board.W; x++ {
		s.winner = s.checkColumnSequence(board, board.State[x], board.StackSize(x))
		if s.winner != Empty {
			return s.winner
		}
	}
	return Empty
}

func (s *Referee) checkHorizontal(board *Board) Player {
	// overlay all columns and determine max stack size
	s.stacksSum = 0
	for x := 0; x < board.W; x++ {
		s.stacksSum |= board.State[x]
	}

	for y := 0; y < 7-bits.LeadingZeros8(uint8(s.stacksSum)); y++ {
		s.lastToken = board.GetCell(0, y)
		s.sameStreak = 1

		for x := 1; x < board.W; x++ {
			s.currentToken = board.GetCell(x, y)
			if s.lastToken == Empty || s.currentToken == Empty || s.currentToken != s.lastToken {
				s.sameStreak = 0
				s.lastToken = s.currentToken
			}
			if s.currentToken != Empty {
				s.sameStreak += 1
			}
			if s.sameStreak >= board.WinStreak {
				return s.currentToken
			}
		}
	}
	return Empty
}

func (s *Referee) checkDiagonals(board *Board) Player {
	// on bottom edge, to right-top
	for xstart := 0; xstart < board.W-board.WinStreak+1; xstart++ {
		s.winner = s.checkDiagonal(board, xstart, 0, +1)
		if s.winner != Empty {
			return s.winner
		}
	}

	// on bottom edge, to left-top
	for xstart := board.WinStreak - 1; xstart < board.W; xstart++ {
		s.winner = s.checkDiagonal(board, xstart, 0, -1)
		if s.winner != Empty {
			return s.winner
		}
	}

	for ystart := 1; ystart < board.H-board.WinStreak+1; ystart++ {
		// on left edge, to right-top
		s.winner = s.checkDiagonal(board, 0, ystart, +1)
		if s.winner != Empty {
			return s.winner
		}

		// on right edge, to left-top
		s.winner = s.checkDiagonal(board, board.W-1, ystart, -1)
		if s.winner != Empty {
			return s.winner
		}
	}

	return Empty
}

func (s *Referee) checkDiagonal(board *Board, xstart, ystart, xstep int) Player {
	s.lastToken = board.GetCell(xstart, ystart)
	s.sameStreak = 1
	x := xstart + xstep
	y := ystart + 1
	for {
		s.currentToken = board.GetCell(x, y)
		if s.lastToken == Empty || s.currentToken == Empty || s.currentToken != s.lastToken {
			s.sameStreak = 0
			s.lastToken = s.currentToken
		}
		if s.currentToken != Empty {
			s.sameStreak += 1
		}
		if s.sameStreak >= board.WinStreak {
			return s.currentToken
		}

		x += xstep
		y++
		if x >= board.W || x < 0 || y >= board.H {
			return Empty
		}
	}
}

// checkColumnSequence checks if there's a winning streak
// input: 						     ooooI1111010001111
// Find 4 consecutive ones
//								>>1: 0ooooI111101000111
//      				   & (>> 1): 0ooooI111000000111
// Find 3 consecutive ones
// 							   >> 1: 00ooooI11110100011
//      				   & (>> 1): 00ooooI11000000011
// Find 2 consecutive ones
// 							   >> 1: 000ooooI1100000001
//      				   & (>> 1): 000ooooI1000000001
// Clear first bits, get last (stack size - (winStreak-1)): & ((1 << (stacksize - (winStreak-1))) -1):
//								  &: 000oooo01111111111
// Is different than 0?
func (s *Referee) checkColumnSequence(board *Board, columnState uint64, stackSize int) Player {
	if stackSize < board.WinStreak {
		return Empty
	}

	onesB := columnState
	onesA := ^columnState

	for i := 0; i < board.WinStreak-1; i++ {
		onesB &= onesB >> 1
		onesA &= onesA >> 1
	}

	var mask uint64 = (1 << (stackSize - (board.WinStreak - 1))) - 1

	if onesA&mask != 0 {
		return PlayerA
	}
	if onesB&mask != 0 {
		return PlayerB
	}
	return Empty
}

func getStackSize(columnState uint64) int {
	return 7 - bits.LeadingZeros8(uint8(columnState))
}

func (s *Referee) whoWonColumn(columnState uint64) Player {
	stackSize := getStackSize(columnState)
	if stackSize < s.winStreak {
		return Empty
	}

	onesB := columnState
	onesA := ^columnState

	for i := 0; i < s.winStreak1; i++ {
		onesB &= onesB >> 1
		onesA &= onesA >> 1
	}

	var mask uint64 = (1 << (stackSize - s.winStreak1)) - 1

	if onesA&mask != 0 {
		return PlayerA
	}
	if onesB&mask != 0 {
		return PlayerB
	}
	return Empty
}

func (s *Referee) hasWonRow(row uint64) bool {
	ones := row
	for i := 0; i < s.winStreak1; i++ {
		ones &= ones >> 1
	}
	return ones != 0
}