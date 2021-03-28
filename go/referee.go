package main

func CheckWinner(board *Board) *Player {
	var winner *Player

	winner = checkVertical(board)
	if winner != nil {
		return winner
	}

	winner = checkHorizontal(board)
	if winner != nil {
		return winner
	}

	winner = checkDiagonals(board)
	if winner != nil {
		return winner
	}

	return nil
}

func checkVertical(board *Board) *Player {
	for _, column := range board.Columns {
		winner := checkContinuousSequence(board, column)
		if winner != nil {
			return winner
		}
	}
	return nil
}

func checkHorizontal(board *Board) *Player {
	maxHeight := 0
	for _, column := range board.Columns {
		if len(column) > maxHeight {
			maxHeight = len(column)
		}
	}

	for y := 0; y < maxHeight; y++ {
		var last *Player = nil
		var e *Player = nil
		streak := 0
		for x := 0; x < board.w; x++ {
			e = board.GetCell(x, y)
			if last == nil || e == nil || *e != *last {
				streak = 0
				last = e
			}
			if e != nil {
				streak += 1
			}
			if streak >= board.winStreak {
				return e
			}
		}
	}
	return nil
}

func checkDiagonals(board *Board) *Player {
	var winner *Player

	// on bottom edge, to right-top
	for xstart := 0; xstart < board.w-board.winStreak+1; xstart++ {
		winner = checkDiagonal(board, xstart, 0, +1)
		if winner != nil {
			return winner
		}
	}

	// on bottom edge, to left-top
	for xstart := board.winStreak - 1; xstart < board.w; xstart++ {
		winner = checkDiagonal(board, xstart, 0, -1)
		if winner != nil {
			return winner
		}
	}

	for ystart := 1; ystart < board.h-board.winStreak+1; ystart++ {
		// on left edge, to right-top
		winner = checkDiagonal(board, 0, ystart, +1)
		if winner != nil {
			return winner
		}

		// on right edge, to left-top
		winner = checkDiagonal(board, board.w-1, ystart, -1)
		if winner != nil {
			return winner
		}
	}

	return nil
}

func checkDiagonal(board *Board, xstart, ystart, xstep int) *Player {
	var last *Player = nil
	var e *Player = nil
	streak := 0
	x := xstart
	y := ystart
	for {
		e = board.GetCell(x, y)
		if last == nil || e == nil || *e != *last {
			streak = 0
			last = e
		}
		if e != nil {
			streak += 1
		}
		if streak >= board.winStreak {
			return e
		}

		x += xstep
		y += 1
		if x >= board.w || x < 0 || y >= board.h {
			return nil
		}
	}
}

// checkSequence checks if there's a winning streak
func CheckSequence(winStreak int, seq []*Player) *Player {
	if len(seq) < winStreak {
		return nil
	}

	var last *Player = nil
	streak := 0
	for _, e := range seq {
		if last == nil || e == nil || *e != *last {
			streak = 0
			last = e
		}
		if e != nil {
			streak += 1
		}
		if streak >= winStreak {
			return e
		}
	}
	return nil
}

func checkContinuousSequence(board *Board, seq []Player) *Player {
	if len(seq) < board.winStreak {
		return nil
	}

	last := seq[0]
	streak := 1
	for i := 1; i < len(seq); i++ {
		e := seq[i]
		if e != last {
			streak = 1
			last = e
			continue
		}
		streak += 1
		if streak >= board.winStreak {
			return &e
		}
	}
	return nil
}
