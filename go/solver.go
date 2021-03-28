package main

type MoveSolver struct {
}

func NewMoveSolver() *MoveSolver {
	return &MoveSolver{}
}

func (s *MoveSolver) MovesEndings(board *Board) []GameEnding {
	return nil
}

func (s *MoveSolver) BestEnding(board *Board) GameEnding {
	return Tie
}

func (s *MoveSolver) BestEndingOnMove(board *Board, player Player, move int) GameEnding {
	return Tie
}
