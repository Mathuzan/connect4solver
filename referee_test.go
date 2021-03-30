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