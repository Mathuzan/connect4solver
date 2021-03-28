package c4solver

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBestResultSimplest4(t *testing.T) {
	board := ParseBoard(`
	....
	ABAB
	ABAB
	ABAB
	`)
	solver := MoveSolver()

	assert.Equal(t, Win,
		solver.BestResultOnMove(grid, PlayerA, 0))
	assert.Equal(t, Lose,
		solver.BestResultOnMove(grid, PlayerA, 1))
	assert.Equal(t, Win,
		solver.BestResultOnMove(grid, PlayerA, 2))
	assert.Equal(t, Lose,
		solver.BestResultOnMove(grid, PlayerA, 3))
	
	assert.Equal(t, Win,
		solver.BestResult(board))
}

func TestBestResultSimpleTie(t *testing.T) {
	board := ParseBoard(`
	..
	AB
	AB
	AB
	`)
	solver := MoveSolver()

	assert.Equal(t, Win,
		solver.BestResultOnMove(grid, PlayerA, 0))
	assert.Equal(t, Tie,
		solver.BestResultOnMove(grid, PlayerA, 1))
	
	assert.Equal(t, Win,
		solver.BestResult(board))
}

func TestBestResult3x3(t *testing.T) {
	board := NewBoard(3, 3, min_win=3)
	solver := MoveSolver()
	results := solver.MovesResults(board)

	assert.Equal(t, []rune{Tie, Tie, Tie},
		results)
}

func TestBestResult3x3Unfair(t *testing.T) {
	board := NewBoard(3, 3, min_win=2)
	solver := MoveSolver()
	results := solver.MovesResults(board)

	assert.Equal(t, Win,
		solver.BestResultOnMove(grid, PlayerA, 1))
	assert.Equal(t, []rune{Win, Win, Win},
		results)
}
