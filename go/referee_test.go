package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNobodyWon(t *testing.T) {
	board := NewBoard(WithSize(7, 6))
	assert.Equal(t,
		nil,
		board.HasWinner())

	board = ParseBoard(`
	. . . . . . .
	. . . . . . .
	. . . . . . B
	. A . . . . B
	. B . B A . A
	. A . A B . B
	`)
	assert.Equal(t, nil, board.HasWinner())
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
	assert.Equal(t, PlayerA, board.HasWinner())
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
	assert.Equal(t, PlayerB, board.HasWinner())
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
	assert.Equal(t, PlayerA, board.HasWinner())

	board = ParseBoard(`
	. . . . . . .
	. . . . . . .
	. B . . . . B
	. A B . . A A
	. B A B A B B
	A B A A B B A
	`)
	assert.Equal(t, PlayerB, board.HasWinner())

	board = ParseBoard(`
	. . . B . . .
	. . B B . . .
	. B B B . . B
	B A A A . A A
	A A A B . B .
	A A B A . B A
	`)
	assert.Equal(t, PlayerB, board.HasWinner())
}

func TestMinStreakCondition(t *testing.T) {
	board := ParseBoard(`
...
.A.
AAB
`, WithWinStreak(2))
	assert.Equal(t, PlayerA, board.HasWinner())
}
