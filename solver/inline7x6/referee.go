package inline7x6

import (
	"github.com/igrek51/connect4solver/solver/common"
)

type Referee struct {
	w         int
	h         int
	winStreak int

	verticalMovesMap            []common.Player
	binaryRowMap                []bool
	diagonalEvaluatorMap        [][]winnerEvaluator
	counterDiagonalEvaluatorMap [][]winnerEvaluator
	winStreak1                  int
}

type winnerEvaluator func(*common.Board, common.Player) bool

type coordinate struct {
	x int
	y int
}

func NewReferee(board *common.Board) *Referee {
	s := &Referee{
		w:          board.W,
		h:          board.H,
		winStreak:  board.WinStreak,
		winStreak1: board.WinStreak - 1,
	}

	// get all possible column layouts and pre-calculate winners for them
	verticalMovesMap := make([]common.Player, 1<<(s.h+1))
	for colState := uint64(0); colState <= (1<<(s.h+1))-1; colState++ {
		verticalMovesMap[colState] = s.whoWonColumn(colState)
	}
	s.verticalMovesMap = verticalMovesMap

	maxSize := s.w
	if s.h > maxSize {
		maxSize = s.h
	}
	// binary row is represented as: 1 - desired token of player, 0 - opponent's token or empty
	binaryRowMap := make([]bool, 1<<maxSize)
	for binaryRow := uint64(0); binaryRow <= (1<<maxSize)-1; binaryRow++ {
		binaryRowMap[binaryRow] = s.hasWonRow(binaryRow)
	}
	s.binaryRowMap = binaryRowMap

	// build functions handling diagonal validation for each coordinate
	s.diagonalEvaluatorMap = make([][]winnerEvaluator, board.W)
	s.counterDiagonalEvaluatorMap = make([][]winnerEvaluator, board.W)
	for x := 0; x < board.W; x++ {
		s.diagonalEvaluatorMap[x] = make([]winnerEvaluator, board.H)
		s.counterDiagonalEvaluatorMap[x] = make([]winnerEvaluator, board.H)
		for y := 0; y < board.H; y++ {
			s.diagonalEvaluatorMap[x][y] = s.buildDiagonalEvaluator(x, y)
			s.counterDiagonalEvaluatorMap[x][y] = s.buildCounterDiagonalEvaluator(x, y)
		}
	}

	return s
}

func (s *Referee) HasPlayerWon(board *common.Board, move int, y int, player common.Player) bool {
	return s.HasPlayerWonVertical(board, move, player) ||
		s.HasPlayerWonHorizontal(board, y, player) ||
		s.HasPlayerWonDiagonal(board, move, y, player)
}

func (s *Referee) HasPlayerWonVertical(board *common.Board, move int, player common.Player) bool {
	return s.verticalMovesMap[board.State[move]] == player
}

func (s *Referee) HasPlayerWonHorizontal(board *common.Board, y int, player common.Player) bool {
	var binaryRow uint64
	if board.GetCell(0, y) == player {
		binaryRow |= 1
	}
	if board.GetCell(1, y) == player {
		binaryRow |= 1 << 1
	}
	if board.GetCell(2, y) == player {
		binaryRow |= 1 << 2
	}
	if board.GetCell(3, y) == player {
		binaryRow |= 1 << 3
	}
	if board.GetCell(4, y) == player {
		binaryRow |= 1 << 4
	}
	if board.GetCell(5, y) == player {
		binaryRow |= 1 << 5
	}
	if board.GetCell(6, y) == player {
		binaryRow |= 1 << 6
	}
	return s.binaryRowMap[binaryRow]
}

func (s *Referee) HasPlayerWonDiagonal(board *common.Board, x int, y int, player common.Player) bool {
	return s.diagonalEvaluatorMap[x][y](board, player) ||
		s.counterDiagonalEvaluatorMap[x][y](board, player)
}

func (s *Referee) HasPlayerWonDiagonalCounter(board *common.Board, startX int, startY int, player common.Player) bool {
	var binaryRow uint64
	y := startY + 3
	for x := startX - 3; x <= startX+3; x++ {
		if x >= 0 && x < 7 && y >= 0 && y < 6 {
			if board.GetCell(x, y) == player {
				binaryRow |= 1 << x
			}
		}
		y--
	}
	return s.binaryRowMap[binaryRow]
}

func getStackSize(columnState uint64) int {
	return common.StackSizeLookup[columnState]
}

func (s *Referee) whoWonColumn(columnState uint64) common.Player {
	stackSize := getStackSize(columnState)
	if stackSize < s.winStreak {
		return common.Empty
	}

	onesB := columnState
	onesA := ^columnState

	for i := 0; i < s.winStreak1; i++ {
		onesB &= onesB >> 1
		onesA &= onesA >> 1
	}

	var mask uint64 = (1 << (stackSize - s.winStreak1)) - 1

	if onesA&mask != 0 {
		return common.PlayerA
	}
	if onesB&mask != 0 {
		return common.PlayerB
	}
	return common.Empty
}

func (s *Referee) hasWonRow(row uint64) bool {
	ones := row
	for i := 0; i < s.winStreak1; i++ {
		ones &= ones >> 1
	}
	return ones != 0
}

func noWinnerEvaluator(*common.Board, common.Player) bool {
	return false
}

func (s *Referee) buildDiagonalEvaluator(startX, startY int) winnerEvaluator {
	coordinates := []coordinate{}
	y := startY - s.winStreak1
	for x := startX - s.winStreak1; x <= startX+s.winStreak1; x++ {
		if x >= 0 && x < s.w && y >= 0 && y < s.h {
			coordinates = append(coordinates, coordinate{x: x, y: y})
		}
		y++
	}

	if len(coordinates) < s.winStreak {
		return noWinnerEvaluator
	}

	return func(board *common.Board, player common.Player) bool {
		var binaryRow uint8
		for i, coordinate := range coordinates {
			if board.GetCell(coordinate.x, coordinate.y) == player {
				binaryRow |= 1 << i
			}
		}
		return s.binaryRowMap[binaryRow]
	}
}

func (s *Referee) buildCounterDiagonalEvaluator(startX, startY int) winnerEvaluator {
	coordinates := []coordinate{}
	y := startY + s.winStreak1
	for x := startX - s.winStreak1; x <= startX+s.winStreak1; x++ {
		if x >= 0 && x < s.w && y >= 0 && y < s.h {
			coordinates = append(coordinates, coordinate{x: x, y: y})
		}
		y--
	}

	if len(coordinates) < s.winStreak {
		return noWinnerEvaluator
	}

	return func(board *common.Board, player common.Player) bool {
		var binaryRow uint8
		for i, coordinate := range coordinates {
			if board.GetCell(coordinate.x, coordinate.y) == player {
				binaryRow |= 1 << i
			}
		}
		return s.binaryRowMap[binaryRow]
	}
}
