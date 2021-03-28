package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyBoardRender(t *testing.T) {
	board := NewBoard(WithSize(7, 6))
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
	board.Throw(1, PlayerA).Throw(1, PlayerB).Throw(3, PlayerA)
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
