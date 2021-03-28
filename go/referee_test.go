package c4solver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNobodyWon(t *testing.T) {
	board := Board(8, 7)
	assert.Equal(t,
		nil,
		board.winner())

	board := ParseBoard(`
	. . . . . . .
	. . . . . . .
	. . . . . . B
	. A . . . . B
	. B . B A . A
	. A . A B . B
	`)
	assert.Equal(t, nil, board.winner())
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
	assert.Equal(t, PlayerA, board.winner())
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
	assert.Equal(t, PlayerB, board.winner())
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
	assert.Equal(t, PlayerA, board.winner())

	board := ParseBoard(`
	. . . . . . .
	. . . . . . .
	. B . . . . B
	. A B . . A A
	. B A B A B B
	A B A A B B A
	`)
	assert.Equal(t, PlayerB, board.winner())

	board := ParseBoard(`
	. . . B . . .
	. . B B . . .
	. B B B . . B
	B A A A . A A
	A A A B . B .
	A A B A . B A
	`)
	assert.Equal(t, PlayerB, board.winner())
}

func TestMinStreakCondition(t *testing.T) {
	board := ParseBoard(`
...
.A.
AAB
`, MinStreakCondition(2))
	assert.Equal(t, PlayerA, board.winner())
}
