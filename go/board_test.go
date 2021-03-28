package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyBoardRender(t *testing.T) {
	board := NewBoard(WithSize(7, 6))
	rendered := board.String()
	assert.Equal(t, `
+---------------+
| . . . . . . . |
| . . . . . . . |
| . . . . . . . |
| . . . . . . . |
| . . . . . . . |
| . . . . . . . |
+---------------+
| 0 1 2 3 4 5 6 |
`,
		rendered,
		"empty board rendered")
}

func TestRenderGridWithTokens(t *testing.T) {
	board := NewBoard(WithSize(7, 6))
	board.Throw(1, PlayerA).Throw(1, PlayerB).Throw(3, PlayerA)
	rendered := board.String()
	assert.Equal(t, `
+---------------+
| . . . . . . . |
| . . . . . . . |
| . . . . . . . |
| . . . . . . . |
| . B . . . . . |
| . A . A . . . |
+---------------+
| 0 1 2 3 4 5 6 |
`,
		rendered,
		"board with thrown tokens rendered")
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
	assert.Equal(t, `
+---------------+
| . . . . . . . |
| . . . . . . . |
| . . . . . . B |
| . . . . . . B |
| . B . . . . A |
| . A . A . . B |
+---------------+
| 0 1 2 3 4 5 6 |
`, rendered)

	board = ParseBoard(`
	. . . . . . .
	. . . . . . .
	. . . . . . B
	. . . . . . B
	. B . . . . A
	. A . A . . B
	`)
	rendered = board.String()
	assert.Equal(t, `
+---------------+
| . . . . . . . |
| . . . . . . . |
| . . . . . . B |
| . . . . . . B |
| . B . . . . A |
| . A . A . . B |
+---------------+
| 0 1 2 3 4 5 6 |
`, rendered)
}
