package solver

import (
	"github.com/igrek51/connect4solver/solver/common"
	"github.com/igrek51/connect4solver/solver/generic_solver"
	"github.com/igrek51/connect4solver/solver/inline7x6"
)

func CreateSolver(board *common.Board) common.IMoveSolver {
	var solver common.IMoveSolver
	// take precedence with inlined optimized solvers
	if board.W == 7 && board.H == 6 {
		solver = inline7x6.NewMoveSolver(board)
	} else {
		solver = generic_solver.NewMoveSolver(board)
	}
	return solver
}
