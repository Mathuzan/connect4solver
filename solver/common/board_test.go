package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyBoardRender(t *testing.T) {
	board := NewBoard(WithSize(7, 6))
	assert.EqualValues(t, 0, board.StackSize(0))
	assert.EqualValues(t, Empty, board.GetCell(0, 0))
	assert.EqualValues(t, Empty, board.GetCell(0, 1))
	rendered := board.String()
	AssertEqualTrimmed(t, rendered, `
+---------------+
| . . . . . . . |
| . . . . . . . |
| . . . . . . . |
| . . . . . . . |
| . . . . . . . |
| . . . . . . . |
+---------------+
| 0 1 2 3 4 5 6 |
`)
}

func TestRenderGridWithTokens(t *testing.T) {
	board := NewBoard(WithSize(7, 6))
	board.Throw(1, PlayerA)
	board.Throw(1, PlayerB)
	board.Throw(3, PlayerA)
	rendered := board.String()
	AssertEqualTrimmed(t, rendered, `
+---------------+
| . . . . . . . |
| . . . . . . . |
| . . . . . . . |
| . . . . . . . |
| . B . . . . . |
| . A . A . . . |
+---------------+
| 0 1 2 3 4 5 6 |
`)
}

func TestParseBoard(t *testing.T) {
	board := ParseBoard(`
.......
.......
......B
......B
.B....A
.A.A..B
`)
	rendered := board.String()
	AssertEqualTrimmed(t, rendered, `
+---------------+
| . . . . . . . |
| . . . . . . . |
| . . . . . . B |
| . . . . . . B |
| . B . . . . A |
| . A . A . . B |
+---------------+
| 0 1 2 3 4 5 6 |
`)

	board = ParseBoard(`
	. . . . . . .
	. . . . . . .
	. . . . . . B
	. . . . . . B
	. B . . . . A
	. A . A . . B
	`)
	rendered = board.String()
	AssertEqualTrimmed(t, rendered, `
+---------------+
| . . . . . . . |
| . . . . . . . |
| . . . . . . B |
| . . . . . . B |
| . B . . . . A |
| . A . A . . B |
+---------------+
| 0 1 2 3 4 5 6 |
`)
}

func TestNextPlayer(t *testing.T) {
	board := ParseBoard(`
. . . . . . .
. . . . . . .
. . . . . . B
. . . A . . B
. B . A . . A
. A . A . . B
`)
	assert.EqualValues(t, PlayerB, board.NextPlayer())
}

func TestBoardOptions(t *testing.T) {
	board := NewBoard(WithSize(4, 3), WithWinStreak(2))
	assert.EqualValues(t, 4, board.W)
	assert.EqualValues(t, 3, board.H)
	assert.EqualValues(t, 2, board.WinStreak)
}

func TestRevert(t *testing.T) {
	board := ParseBoard(`
.......
.......
......B
......B
.B....A
.A.A..B
`)
	board.Revert(1, 1)
	board.Revert(3, 0)
	board.Revert(6, 3)
	board.Revert(6, 2)
	board.Revert(6, 1)
	rendered := board.String()
	AssertEqualTrimmed(t, rendered, `
+---------------+
| . . . . . . . |
| . . . . . . . |
| . . . . . . . |
| . . . . . . . |
| . . . . . . . |
| . A . . . . B |
+---------------+
| 0 1 2 3 4 5 6 |
`)
}

func TestApplyMoves(t *testing.T) {
	board := NewBoard(WithSize(7, 6))
	board.ApplyMoves("0035666")
	rendered := board.String()
	AssertEqualTrimmed(t, rendered, `
+---------------+
| . . . . . . . |
| . . . . . . . |
| . . . . . . . |
| . . . . . . A |
| B . . . . . B |
| A . . A . B A |
+---------------+
| 0 1 2 3 4 5 6 |
`)
}
