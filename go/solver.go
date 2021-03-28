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
		endings[move] = ending
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
	boardAfter := board.Clone().Throw(move, player)
	nextPlayer := oppositePlayer(player)

	winner := board.HasWinner()
	if winner != nil {
		if *winner == PlayerA {
			return Win
		} else {
			return Lose
		}
	}

	ending := s.NextMoveEnding(boardAfter, nextPlayer)
	return ending
}

func (s *MoveSolver) NextMoveEnding(board *Board, player Player) GameEnding {
	// find further possible moves
	possibleMovesEndings := []GameEnding{}
	for move := 0; move < board.w; move++ {
		if board.CanMakeMove(move) {
			moveEnding := s.BestEndingOnMove(board, player, move)
			if player == PlayerA && moveEnding == Win {
				return Win
			}
			if player == PlayerB && moveEnding == Lose {
				return Lose
			}

			possibleMovesEndings = append(possibleMovesEndings, moveEnding)
		}
	}

	if len(possibleMovesEndings) == 0 {
		return Tie
	}

	if player == PlayerA {
		return maxPossibleMove(possibleMovesEndings)
	} else {
		return minPossibleMove(possibleMovesEndings)
	}
}

func oppositePlayer(player Player) Player {
	if player == PlayerA {
		return PlayerB
	} else {
		return PlayerA
	}
}

func maxPossibleMove(endings []GameEnding) GameEnding {
	maxr := endings[0]
	for _, ending := range endings {
		if MoveResultsWeights[ending] > MoveResultsWeights[maxr] {
			maxr = ending
		}
	}
	return maxr
}

func minPossibleMove(endings []GameEnding) GameEnding {
	minr := endings[0]
	for _, ending := range endings {
		if MoveResultsWeights[ending] < MoveResultsWeights[minr] {
			minr = ending
		}
	}
	return minr
}
