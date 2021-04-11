package solver

import (
	"fmt"
	"testing"

	. "github.com/igrek51/connect4solver/solver/common"
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

func TestCachedResultsCount(t *testing.T) {
	board := NewBoard(WithSize(3, 2), WithWinStreak(3))
	solver := NewMoveSolver(board)

	endings := solver.MovesEndings(board)
	assert.Equal(t, []Player{Empty, Empty, Empty}, endings)

	depth := 1
	fmt.Printf("cached boards for depth: %d, len: %d\n", depth, len(solver.cache.depthCaches[depth]))
	for key := range solver.cache.depthCaches[depth] {
		var state BoardKey
		state[0] = uint64(uint8(key))
		state[1] = uint64(uint8(key >> 8))
		state[2] = uint64(uint8(key >> 16))
		board.State = state
		fmt.Println(board.String())
	}
	assert.Equal(t, 2, len(solver.cache.depthCaches[0]))
	assert.Equal(t, 5, len(solver.cache.depthCaches[1]))
}

func Test7x6Solver(t *testing.T) {
	board := ParseBoard(`
	.......
	.......
	.......
	ABABABA
	ABABABB
	ABABABA
	`)
	solver := CreateSolver(board)
	endings := solver.MovesEndings(board)
	assert.Equal(t, []Player{PlayerA, PlayerB, PlayerA, PlayerB, PlayerA, PlayerB, PlayerA}, endings)
}

func BenchmarkMoveSolver4x4(b *testing.B) {
	board := NewBoard(WithSize(4, 4), WithWinStreak(4))
	b.ResetTimer()
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		solver := CreateSolver(board)
		b.StartTimer()
		solver.MovesEndings(board)
		b.StopTimer()
	}
}
