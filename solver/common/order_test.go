package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMovesOrderEvenBoard(t *testing.T) {
	board := NewBoard(WithSize(6, 5), WithWinStreak(4))
	movesOrder := CalculateMovesOrder(board)
	assert.Equal(t, []int{3, 2, 4, 1, 5, 0}, movesOrder)
}

func TestMovesOrderOddBoard(t *testing.T) {
	board := NewBoard(WithSize(7, 6), WithWinStreak(4))
	movesOrder := CalculateMovesOrder(board)
	assert.Equal(t, []int{3, 4, 2, 5, 1, 6, 0}, movesOrder)
}
