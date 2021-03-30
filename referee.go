package main

import "math/bits"

var winner Player
var lastToken uint8
var lastTokenPlayer Player
var currentToken Player
var sameStreak int
var stacksSum uint8

func CheckWinner(board *Board) Player {
	winner = checkVertical(board)
	if winner != Empty {
		return winner
	}

	winner = checkHorizontal(board)
	if winner != Empty {
		return winner
	}

	winner = checkDiagonals(board)
	if winner != Empty {
		return winner
	}

	return Empty
}

func checkVertical(board *Board) Player {
	for x := 0; x < board.w; x++ {
		winner = checkColumnSequence(board, board.state[x], board.stackSize(x))
		if winner != Empty {
			return winner
		}
	}
	return Empty
}

func checkHorizontal(board *Board) Player {
	// overlay all columns and determine max stack size
	stacksSum = 0
	for x := 0; x < board.w; x++ {
		stacksSum |= board.state[x]
	}

	for y := 0; y < 7-bits.LeadingZeros8(stacksSum); y++ {
		lastTokenPlayer = board.GetCell(0, y)
		sameStreak = 1

		for x := 1; x < board.w; x++ {
			currentToken = board.GetCell(x, y)
			if lastTokenPlayer == Empty || currentToken == Empty || currentToken != lastTokenPlayer {
				sameStreak = 0
				lastTokenPlayer = currentToken
			}
			if currentToken != Empty {
				sameStreak += 1
			}
			if sameStreak >= board.winStreak {
				return currentToken
			}
		}
	}
	return Empty
}

func checkDiagonals(board *Board) Player {
	// on bottom edge, to right-top
	for xstart := 0; xstart < board.w-board.winStreak+1; xstart++ {
		winner = checkDiagonal(board, xstart, 0, +1)
		if winner != Empty {
			return winner
		}
	}

	// on bottom edge, to left-top
	for xstart := board.winStreak - 1; xstart < board.w; xstart++ {
		winner = checkDiagonal(board, xstart, 0, -1)
		if winner != Empty {
			return winner
		}
	}

	for ystart := 1; ystart < board.h-board.winStreak+1; ystart++ {
		// on left edge, to right-top
		winner = checkDiagonal(board, 0, ystart, +1)
		if winner != Empty {
			return winner
		}

		// on right edge, to left-top
		winner = checkDiagonal(board, board.w-1, ystart, -1)
		if winner != Empty {
			return winner
		}
	}

	return Empty
}

func checkDiagonal(board *Board, xstart, ystart, xstep int) Player {
	lastTokenPlayer = board.GetCell(xstart, ystart)
	sameStreak = 1
	x := xstart + xstep
	y := ystart + 1
	for {
		currentToken = board.GetCell(x, y)
		if lastTokenPlayer == Empty || currentToken == Empty || currentToken != lastTokenPlayer {
			sameStreak = 0
			lastTokenPlayer = currentToken
		}
		if currentToken != Empty {
			sameStreak += 1
		}
		if sameStreak >= board.winStreak {
			return currentToken
		}

		x += xstep
		y++
		if x >= board.w || x < 0 || y >= board.h {
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
func checkColumnSequence(board *Board, columnState uint8, stackSize int) Player {
	if stackSize < board.winStreak {
		return Empty
	}

	onesB := columnState
	onesA := ^columnState

	for i := 0; i < board.winStreak-1; i++ {
		onesB &= onesB >> 1
		onesA &= onesA >> 1
	}

	var mask uint8 = (1 << (stackSize - (board.winStreak - 1))) - 1

	if onesA&mask != 0 {
		return PlayerA
	}
	if onesB&mask != 0 {
		return PlayerB
	}
	return Empty
}
