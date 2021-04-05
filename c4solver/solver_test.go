package c4solver

import (
	"testing"

	. "github.com/igrek51/connect4solver/c4solver/common"
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

	assert.Equal(t, PlayerA,
		solver.BestEndingOnMove(board, PlayerA, 0))
	assert.Equal(t, PlayerB,
		solver.BestEndingOnMove(board, PlayerA, 1))
	assert.Equal(t, PlayerA,
		solver.BestEndingOnMove(board, PlayerA, 2))
	assert.Equal(t, PlayerB,
		solver.BestEndingOnMove(board, PlayerA, 3))
}

func TestBestResultSimpleTie(t *testing.T) {
	board := ParseBoard(`
	..
	AB
	AB
	AB
	`)
	solver := NewMoveSolver(board)

	assert.Equal(t, PlayerA,
		solver.BestEndingOnMove(board, PlayerA, 0))
	assert.Equal(t, Empty,
		solver.BestEndingOnMove(board, PlayerA, 1))
}

func TestBestResult3x3(t *testing.T) {
	board := NewBoard(WithSize(3, 3), WithWinStreak(3))
	solver := NewMoveSolver(board)
	endings := solver.MovesEndings(board)

	assert.Equal(t, []Player{Empty, Empty, Empty},
		endings)
}

func TestBestResult3x3Unfair(t *testing.T) {
	board := NewBoard(WithSize(3, 3), WithWinStreak(2))
	solver := NewMoveSolver(board)
	endings := solver.MovesEndings(board)

	assert.Equal(t, PlayerA,
		solver.BestEndingOnMove(board, PlayerA, 1))
	assert.Equal(t, []Player{PlayerA, PlayerA, PlayerA},
		endings)
}

func TestBestResult2x2Unfair(t *testing.T) {
	board := NewBoard(WithSize(2, 2), WithWinStreak(2))
	solver := NewMoveSolver(board)
	endings := solver.MovesEndings(board)

	assert.Equal(t, []Player{PlayerA, PlayerA},
		endings)
}

func BenchmarkMoveSolver4x4(b *testing.B) {
	board := NewBoard(WithSize(4, 4), WithWinStreak(4))
	solver := NewMoveSolver(board)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		solver.MovesEndings(board)
	}
}
