package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNobodyWon(t *testing.T) {
	board := NewBoard(WithSize(7, 6))
	referee := NewReferee(board)
	assert.EqualValues(t, Empty, referee.HasWinner(board))

	board = ParseBoard(`
	. . . . . . .
	. . . . . . .
	. . . . . . B
	. A . . . . B
	. B . B A . A
	. A . A B . B
	`)
	referee = NewReferee(board)
	assert.EqualValues(t, Empty, referee.HasWinner(board))
	assert.EqualValues(t, false, referee.HasPlayerWon(board, 1, 2, PlayerA))
	assert.EqualValues(t, false, referee.HasPlayerWon(board, 6, 3, PlayerB))
}

func TestWonVertical(t *testing.T) {
	board := ParseBoard(`
	. . . . . . .
	. . . . . . .
	. . . . . . A
	. . . . . . A
	. B . . . . A
	. A . A . . A
	`)
	referee := NewReferee(board)
	assert.EqualValues(t, PlayerA, referee.HasWinner(board))
	assert.EqualValues(t, true, referee.HasPlayerWon(board, 6, 3, PlayerA))
}

func TestWonHorizontal(t *testing.T) {
	board := ParseBoard(`
	. . . . . . .
	. . . . . . .
	. . . . . . A
	. . . . . . A
	. B . B B B B
	A A . A B B A
	`)
	referee := NewReferee(board)
	winner := referee.HasWinner(board)
	assert.EqualValues(t, PlayerB, winner)
	assert.EqualValues(t, true, referee.HasPlayerWon(board, 3, 1, PlayerB))
	assert.EqualValues(t, true, referee.HasPlayerWon(board, 4, 1, PlayerB))
	assert.EqualValues(t, true, referee.HasPlayerWon(board, 5, 1, PlayerB))
	assert.EqualValues(t, true, referee.HasPlayerWon(board, 6, 1, PlayerB))

	assert.EqualValues(t, false, referee.HasPlayerWon(board, 6, 2, PlayerA))
}

func TestWonDiagonal(t *testing.T) {
	board := ParseBoard(`
	. . . . . . .
	. . . . . . .
	. . . . . . A
	. . . . . A A
	. B . B A B B
	A A . A B B A
	`)
	referee := NewReferee(board)
	assert.EqualValues(t, PlayerA, referee.HasWinner(board))
	assert.EqualValues(t, true, referee.HasPlayerWon(board, 3, 0, PlayerA))
	assert.EqualValues(t, true, referee.HasPlayerWon(board, 4, 1, PlayerA))
	assert.EqualValues(t, true, referee.HasPlayerWon(board, 5, 2, PlayerA))
	assert.EqualValues(t, true, referee.HasPlayerWon(board, 6, 3, PlayerA))

	board = ParseBoard(`
	. . . . . . .
	. . . . . . .
	. B . . . . B
	. A B . . A A
	. B A B A B B
	A B A A B B A
	`)
	assert.EqualValues(t, PlayerB, referee.HasWinner(board))
	assert.EqualValues(t, true, referee.HasPlayerWon(board, 1, 3, PlayerB))
	assert.EqualValues(t, true, referee.HasPlayerWon(board, 2, 2, PlayerB))
	assert.EqualValues(t, true, referee.HasPlayerWon(board, 3, 1, PlayerB))
	assert.EqualValues(t, true, referee.HasPlayerWon(board, 4, 0, PlayerB))

	board = ParseBoard(`
	. . . B . . .
	. . B B . . .
	. B B B . . B
	B A A A . A A
	A A A B . B .
	A A B A . B A
	`)
	assert.EqualValues(t, PlayerB, referee.HasWinner(board))
	assert.EqualValues(t, true, referee.HasPlayerWon(board, 0, 2, PlayerB))
	assert.EqualValues(t, true, referee.HasPlayerWon(board, 1, 3, PlayerB))
	assert.EqualValues(t, true, referee.HasPlayerWon(board, 2, 4, PlayerB))
	assert.EqualValues(t, true, referee.HasPlayerWon(board, 3, 5, PlayerB))
}

func TestMinStreakCondition(t *testing.T) {
	board := ParseBoard(`
...
.A.
AAB
`, WithWinStreak(2))
	referee := NewReferee(board)
	assert.EqualValues(t, PlayerA, referee.HasWinner(board))
	assert.EqualValues(t, true, referee.HasPlayerWon(board, 0, 0, PlayerA))
	assert.EqualValues(t, true, referee.HasPlayerWon(board, 1, 0, PlayerA))
	assert.EqualValues(t, true, referee.HasPlayerWon(board, 1, 1, PlayerA))
}

func TestPlayerWonVerticalOnMove(t *testing.T) {
	board := ParseBoard(".\n.\n.\n.\n.\n.")
	referee := NewReferee(board)

	assert.EqualValues(t, false, referee.HasPlayerWonVertical(board, 0, PlayerA))
	assert.EqualValues(t, false, referee.HasPlayerWonVertical(board, 0, PlayerB))

	board = ParseBoard(".\n.\n.\n.\n.\nA")
	assert.EqualValues(t, false, referee.HasPlayerWonVertical(board, 0, PlayerA))
	assert.EqualValues(t, false, referee.HasPlayerWonVertical(board, 0, PlayerB))
	board = ParseBoard(".\n.\n.\n.\nB\nA")
	assert.EqualValues(t, false, referee.HasPlayerWonVertical(board, 0, PlayerA))
	assert.EqualValues(t, false, referee.HasPlayerWonVertical(board, 0, PlayerB))
	board = ParseBoard(".\n.\n.\nA\nA\nA")
	assert.EqualValues(t, false, referee.HasPlayerWonVertical(board, 0, PlayerA))
	assert.EqualValues(t, false, referee.HasPlayerWonVertical(board, 0, PlayerB))
	board = ParseBoard(".\nA\nB\nB\nA\nA")
	assert.EqualValues(t, false, referee.HasPlayerWonVertical(board, 0, PlayerA))
	assert.EqualValues(t, false, referee.HasPlayerWonVertical(board, 0, PlayerB))
	board = ParseBoard("B\nB\nB\nA\nA\nA")
	assert.EqualValues(t, false, referee.HasPlayerWonVertical(board, 0, PlayerA))
	assert.EqualValues(t, false, referee.HasPlayerWonVertical(board, 0, PlayerB))

	board = ParseBoard("B\nB\nB\nB\nB\nB")
	assert.EqualValues(t, false, referee.HasPlayerWonVertical(board, 0, PlayerA))
	assert.EqualValues(t, true, referee.HasPlayerWonVertical(board, 0, PlayerB))
	board = ParseBoard("A\nA\nA\nA\nA\nA")
	assert.EqualValues(t, PlayerA, referee.checkVertical(board))
	assert.EqualValues(t, true, referee.HasPlayerWonVertical(board, 0, PlayerA))
	assert.EqualValues(t, false, referee.HasPlayerWonVertical(board, 0, PlayerB))

	board = ParseBoard(".\n.\nA\nA\nA\nA")
	assert.EqualValues(t, true, referee.HasPlayerWonVertical(board, 0, PlayerA))
	assert.EqualValues(t, false, referee.HasPlayerWonVertical(board, 0, PlayerB))
	board = ParseBoard(".\n.\nB\nB\nB\nB")
	assert.EqualValues(t, false, referee.HasPlayerWonVertical(board, 0, PlayerA))
	assert.EqualValues(t, true, referee.HasPlayerWonVertical(board, 0, PlayerB))

	board = ParseBoard(".\nB\nB\nB\nB\nA")
	assert.EqualValues(t, false, referee.HasPlayerWonVertical(board, 0, PlayerA))
	assert.EqualValues(t, true, referee.HasPlayerWonVertical(board, 0, PlayerB))
	board = ParseBoard("B\nB\nB\nB\nA\nA")
	assert.EqualValues(t, false, referee.HasPlayerWonVertical(board, 0, PlayerA))
	assert.EqualValues(t, true, referee.HasPlayerWonVertical(board, 0, PlayerB))

	board = ParseBoard("A\nB\nB\nB\nB\nA")
	assert.EqualValues(t, false, referee.HasPlayerWonVertical(board, 0, PlayerA))
	assert.EqualValues(t, true, referee.HasPlayerWonVertical(board, 0, PlayerB))

	board = ParseBoard("A\nA\nA\nA\nB\nA")
	assert.EqualValues(t, true, referee.HasPlayerWonVertical(board, 0, PlayerA))
	assert.EqualValues(t, false, referee.HasPlayerWonVertical(board, 0, PlayerB))
}
