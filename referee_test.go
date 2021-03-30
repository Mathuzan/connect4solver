package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNobodyWon(t *testing.T) {
	board := NewBoard(WithSize(7, 6))
	assert.EqualValues(t, Empty, board.HasWinner())

	board = ParseBoard(`
	. . . . . . .
	. . . . . . .
	. . . . . . B
	. A . . . . B
	. B . B A . A
	. A . A B . B
	`)
	assert.EqualValues(t, Empty, board.HasWinner())
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
	assert.EqualValues(t, PlayerA, board.HasWinner())
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
	winner := board.HasWinner()
	assert.EqualValues(t, PlayerB, winner)
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
	assert.EqualValues(t, PlayerA, board.HasWinner())

	board = ParseBoard(`
	. . . . . . .
	. . . . . . .
	. B . . . . B
	. A B . . A A
	. B A B A B B
	A B A A B B A
	`)
	assert.EqualValues(t, PlayerB, board.HasWinner())

	board = ParseBoard(`
	. . . B . . .
	. . B B . . .
	. B B B . . B
	B A A A . A A
	A A A B . B .
	A A B A . B A
	`)
	assert.EqualValues(t, PlayerB, board.HasWinner())
}

func TestMinStreakCondition(t *testing.T) {
	board := ParseBoard(`
...
.A.
AAB
`, WithWinStreak(2))
	assert.EqualValues(t, PlayerA, board.HasWinner())
}

func TestCheckVertical(t *testing.T) {
	board := ParseBoard(".\n.\n.\n.\n.\n.")
	assert.EqualValues(t, Empty, checkVertical(board))

	board = ParseBoard(".\n.\n.\n.\n.\nA")
	assert.EqualValues(t, Empty, checkVertical(board))
	board = ParseBoard(".\n.\n.\n.\nB\nA")
	assert.EqualValues(t, Empty, checkVertical(board))
	board = ParseBoard(".\n.\n.\nA\nA\nA")
	assert.EqualValues(t, Empty, checkVertical(board))
	board = ParseBoard(".\nA\nB\nB\nA\nA")
	assert.EqualValues(t, Empty, checkVertical(board))
	board = ParseBoard("B\nB\nB\nA\nA\nA")
	assert.EqualValues(t, Empty, checkVertical(board))

	board = ParseBoard("B\nB\nB\nB\nB\nB")
	assert.EqualValues(t, PlayerB, checkVertical(board))
	board = ParseBoard("A\nA\nA\nA\nA\nA")
	assert.EqualValues(t, PlayerA, checkVertical(board))

	board = ParseBoard(".\n.\nA\nA\nA\nA")
	assert.EqualValues(t, PlayerA, checkVertical(board))
	board = ParseBoard(".\n.\nB\nB\nB\nB")
	assert.EqualValues(t, PlayerB, checkVertical(board))

	board = ParseBoard(".\nB\nB\nB\nB\nA")
	assert.EqualValues(t, PlayerB, checkVertical(board))
	board = ParseBoard("B\nB\nB\nB\nA\nA")
	assert.EqualValues(t, PlayerB, checkVertical(board))

	board = ParseBoard("A\nB\nB\nB\nB\nA")
	assert.EqualValues(t, PlayerB, checkVertical(board))

	board = ParseBoard("A\nA\nA\nA\nB\nA")
	assert.EqualValues(t, PlayerA, checkVertical(board))
}
