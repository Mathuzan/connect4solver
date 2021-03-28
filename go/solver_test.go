package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBestResultSimplest4(t *testing.T) {
	board := ParseBoard(`
	....
	ABAB
	ABAB
	ABAB
	`)
	solver := NewMoveSolver()

	assert.Equal(t, Win,
		solver.BestEndingOnMove(board, PlayerA, 0))
	assert.Equal(t, Lose,
		solver.BestEndingOnMove(board, PlayerA, 1))
	assert.Equal(t, Win,
		solver.BestEndingOnMove(board, PlayerA, 2))
	assert.Equal(t, Lose,
		solver.BestEndingOnMove(board, PlayerA, 3))

	assert.Equal(t, Win,
		solver.BestEnding(board))
}

func TestBestResultSimpleTie(t *testing.T) {
	board := ParseBoard(`
	..
	AB
	AB
	AB
	`)
	solver := NewMoveSolver()

	assert.Equal(t, Win,
		solver.BestEndingOnMove(board, PlayerA, 0))
	assert.Equal(t, Tie,
		solver.BestEndingOnMove(board, PlayerA, 1))

	assert.Equal(t, Win,
		solver.BestEnding(board))
}

func TestBestResult3x3(t *testing.T) {
	board := NewBoard(WithSize(3, 3), WithWinStreak(3))
	solver := NewMoveSolver()
	endings := solver.MovesEndings(board)

	assert.Equal(t, []GameEnding{Tie, Tie, Tie},
		endings)
}

func TestBestResult3x3Unfair(t *testing.T) {
	board := NewBoard(WithSize(3, 3), WithWinStreak(2))
	solver := NewMoveSolver()
	endings := solver.MovesEndings(board)

	assert.Equal(t, Win,
		solver.BestEndingOnMove(board, PlayerA, 1))
	assert.Equal(t, []GameEnding{Win, Win, Win},
		endings)
}
