package solver

import (
	"testing"

	. "github.com/igrek51/connect4solver/solver/common"
	"github.com/stretchr/testify/assert"
)

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
