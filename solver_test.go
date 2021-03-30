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
	solver := NewMoveSolver(board)

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
	solver := NewMoveSolver(board)

	assert.Equal(t, Win,
		solver.BestEndingOnMove(board, PlayerA, 0))
	assert.Equal(t, Tie,
		solver.BestEndingOnMove(board, PlayerA, 1))

	assert.Equal(t, Win,
		solver.BestEnding(board))
}

func TestBestResult3x3(t *testing.T) {
	board := NewBoard(WithSize(3, 3), WithWinStreak(3))
	solver := NewMoveSolver(board)
	endings := solver.MovesEndings(board)

	assert.Equal(t, []GameEnding{Tie, Tie, Tie},
		endings)
}

func TestBestResult3x3Unfair(t *testing.T) {
	board := NewBoard(WithSize(3, 3), WithWinStreak(2))
	solver := NewMoveSolver(board)
	endings := solver.MovesEndings(board)

	assert.Equal(t, Win,
		solver.BestEndingOnMove(board, PlayerA, 1))
	assert.Equal(t, []GameEnding{Win, Win, Win},
		endings)
}

func TestBestResult2x2Unfair(t *testing.T) {
	board := NewBoard(WithSize(2, 2), WithWinStreak(2))
	solver := NewMoveSolver(board)
	endings := solver.MovesEndings(board)

	assert.Equal(t, []GameEnding{Win, Win},
		endings)
}

func BenchmarkMoveSolver3x3(b *testing.B) {
	board := NewBoard(WithSize(3, 3), WithWinStreak(3))
	solver := NewMoveSolver(board)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		solver.BestEnding(board)
	}
}