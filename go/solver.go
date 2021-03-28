package main

type MoveSolver struct {
}

func NewMoveSolver() *MoveSolver {
	return &MoveSolver{}
}

func (s *MoveSolver) MovesEndings(board *Board) []GameEnding {
	endings := make([]GameEnding, board.w)
	player := board.NextPlayer()
	for move := 0; move < board.w; move++ {
		ending := s.BestEndingOnMove(board, player, move)
		endings = append(endings, ending)
	}

	return endings
}

func (s *MoveSolver) BestEnding(board *Board) GameEnding {
	endings := s.MovesEndings(board)
	bestEnding := Lose
	for _, ending := range endings {
		if MoveResultsWeights[ending] > MoveResultsWeights[bestEnding] {
			bestEnding = ending
		}
	}
	return bestEnding
}

func (s *MoveSolver) BestEndingOnMove(board *Board, player Player, move int) GameEnding {
	return Tie
}
