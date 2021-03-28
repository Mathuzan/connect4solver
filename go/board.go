package c4solver

type Board struct {
	w       int
	h       int
	min_win int
	columns [][]int
}

func NewBoard(width int, height int, min_win int) *Board {
	columns := [][]int{}
	for i := 0; i < width; i++ {
		columns = append(columns, []int{})
	}

	return &Board{
		w:       width,
		h:       height,
		min_win: min_win,
		columns: columns,
	}
}
