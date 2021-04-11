package common

func CalculateMovesOrder(board *Board) []int {
	movesOrder := []int{}
	for x := 0; x < board.W; x++ {
		var move int
		if x%2 == 0 {
			move = x / 2
		} else {
			move = board.W - 1 - x/2
		}
		movesOrder = append([]int{move}, movesOrder...)
	}
	return movesOrder
}
