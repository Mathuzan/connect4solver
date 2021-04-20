package generic_solver

import (
	"github.com/igrek51/connect4solver/solver/common"
)

func (s *Referee) HasWinner(board *common.Board) common.Player {
	s.winner = s.checkVertical(board)
	if s.winner != common.Empty {
		return s.winner
	}

	s.winner = s.checkHorizontal(board)
	if s.winner != common.Empty {
		return s.winner
	}

	s.winner = s.checkDiagonals(board)
	if s.winner != common.Empty {
		return s.winner
	}

	return common.Empty
}

func (s *Referee) checkVertical(board *common.Board) common.Player {
	for x := 0; x < board.W; x++ {
		s.winner = s.checkColumnSequence(board, board.State[x], board.StackSize(x))
		if s.winner != common.Empty {
			return s.winner
		}
	}
	return common.Empty
}

func (s *Referee) checkHorizontal(board *common.Board) common.Player {
	// overlay all columns and determine max stack size
	s.stacksSum = 0
	for x := 0; x < board.W; x++ {
		s.stacksSum |= board.State[x]
	}

	for y := 0; y < common.StackSizeLookup[s.stacksSum]; y++ {
		s.lastToken = board.GetCell(0, y)
		s.sameStreak = 1

		for x := 1; x < board.W; x++ {
			s.currentToken = board.GetCell(x, y)
			if s.lastToken == common.Empty || s.currentToken == common.Empty || s.currentToken != s.lastToken {
				s.sameStreak = 0
				s.lastToken = s.currentToken
			}
			if s.currentToken != common.Empty {
				s.sameStreak += 1
			}
			if s.sameStreak >= board.WinStreak {
				return s.currentToken
			}
		}
	}
	return common.Empty
}

func (s *Referee) checkDiagonals(board *common.Board) common.Player {
	// on bottom edge, to right-top
	for xstart := 0; xstart < board.W-board.WinStreak+1; xstart++ {
		s.winner = s.checkDiagonal(board, xstart, 0, +1)
		if s.winner != common.Empty {
			return s.winner
		}
	}

	// on bottom edge, to left-top
	for xstart := board.WinStreak - 1; xstart < board.W; xstart++ {
		s.winner = s.checkDiagonal(board, xstart, 0, -1)
		if s.winner != common.Empty {
			return s.winner
		}
	}

	for ystart := 1; ystart < board.H-board.WinStreak+1; ystart++ {
		// on left edge, to right-top
		s.winner = s.checkDiagonal(board, 0, ystart, +1)
		if s.winner != common.Empty {
			return s.winner
		}

		// on right edge, to left-top
		s.winner = s.checkDiagonal(board, board.W-1, ystart, -1)
		if s.winner != common.Empty {
			return s.winner
		}
	}

	return common.Empty
}

func (s *Referee) checkDiagonal(board *common.Board, xstart, ystart, xstep int) common.Player {
	s.lastToken = board.GetCell(xstart, ystart)
	s.sameStreak = 1
	x := xstart + xstep
	y := ystart + 1
	for {
		s.currentToken = board.GetCell(x, y)
		if s.lastToken == common.Empty || s.currentToken == common.Empty || s.currentToken != s.lastToken {
			s.sameStreak = 0
			s.lastToken = s.currentToken
		}
		if s.currentToken != common.Empty {
			s.sameStreak += 1
		}
		if s.sameStreak >= board.WinStreak {
			return s.currentToken
		}

		x += xstep
		y++
		if x >= board.W || x < 0 || y >= board.H {
			return common.Empty
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
func (s *Referee) checkColumnSequence(board *common.Board, columnState uint64, stackSize int) common.Player {
	if stackSize < board.WinStreak {
		return common.Empty
	}

	onesB := columnState
	onesA := ^columnState

	for i := 0; i < board.WinStreak-1; i++ {
		onesB &= onesB >> 1
		onesA &= onesA >> 1
	}

	var mask uint64 = (1 << (stackSize - (board.WinStreak - 1))) - 1

	if onesA&mask != 0 {
		return common.PlayerA
	}
	if onesB&mask != 0 {
		return common.PlayerB
	}
	return common.Empty
}
