package main

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type Board struct {
	w           int
	h           int
	winStreak   int
	Columns     [][]*Player
	ColumnSizes []int
}

func NewBoard(options ...Option) *Board {
	// set defaults
	b := &Board{
		w:         4,
		h:         4,
		winStreak: 4,
	}

	// apply options
	for _, opt := range options {
		err := opt(b)
		if err != nil {
			panic(errors.Wrap(err, "applying option"))
		}
	}

	// validate & post-process
	b.Columns = make([][]*Player, b.w)
	b.ColumnSizes = make([]int, b.w)
	for x := 0; x < b.w; x++ {
		b.Columns[x] = make([]*Player, b.h)
		b.ColumnSizes[x] = 0
	}

	return b
}

type Option func(*Board) error

func WithSize(width int, height int) Option {
	return func(b *Board) error {
		b.w = width
		b.h = height
		return nil
	}
}

func WithWinStreak(winStreak int) Option {
	return func(b *Board) error {
		b.winStreak = winStreak
		return nil
	}
}

// GetCell return token at given cell coordinates (axes oriented top-right)
func (b *Board) GetCell(x int, y int) *Player {
	column := b.Columns[x]
	if y >= b.ColumnSizes[x] {
		return nil
	}
	return column[y]
}

func (b *Board) Throw(x int, player Player) *Board {
	if b.ColumnSizes[x] >= b.h {
		panic("column is already full")
	}
	b.Columns[x][b.ColumnSizes[x]] = &player
	b.ColumnSizes[x]++
	return b
}

func (b *Board) Revert(x int) {
	if b.ColumnSizes[x] <= 0 {
		panic("cant pull out from empty column")
	}
	b.ColumnSizes[x]--
}

const (
	ColouredPlayerA = "\u001b[33;1mA\u001b[0m"
	ColouredPlayerB = "\u001b[31;1mB\u001b[0m"
)

func (b *Board) String() string {
	lines := []string{}
	var line string

	line = "+-"
	for i := 0; i < b.w; i++ {
		line += "--"
	}
	line += "+"
	lines = append(lines, line)

	for y := b.h - 1; y >= 0; y-- {
		rowCells := []string{}
		for x := 0; x < b.w; x++ {
			cell := b.GetCell(x, y)
			if cell == nil {
				rowCells = append(rowCells, string(CellEmpty))
			} else {
				if *cell == PlayerA {
					rowCells = append(rowCells, ColouredPlayerA)
				} else if *cell == PlayerB {
					rowCells = append(rowCells, ColouredPlayerB)
				}
			}
		}
		line = "| " + strings.Join(rowCells, " ") + " |"
		lines = append(lines, line)
	}

	line = "+-"
	for i := 0; i < b.w; i++ {
		line += "--"
	}
	line += "+"
	lines = append(lines, line)

	coordinates := []string{}
	for x := 0; x < b.w; x++ {
		coordinates = append(coordinates, fmt.Sprint(x))
	}
	line = "| " + strings.Join(coordinates, " ") + " |"
	lines = append(lines, line)

	return strings.Join(lines, "\n")
}

func (b *Board) HasWinner() *Player {
	return CheckWinner(b)
}

func (b *Board) NextPlayer() Player {
	tokensA := 0
	tokensB := 0
	for x, column := range b.Columns {
		for y := 0; y < b.ColumnSizes[x]; y++ {
			cell := column[y]
			if *cell == PlayerA {
				tokensA++
			} else if *cell == PlayerB {
				tokensB++
			}
		}
	}
	if tokensA > tokensB {
		return PlayerB
	} else {
		return PlayerA
	}
}

func (b *Board) CanMakeMove(x int) bool {
	return b.ColumnSizes[x] < b.h
}

func (b *Board) Clone() *Board {
	columns := make([][]*Player, b.w)
	columnSizes := make([]int, b.w)
	for x := 0; x < b.w; x++ {
		columnSizes[x] = b.ColumnSizes[x]
		columns[x] = make([]*Player, b.h)
		copy(columns[x], b.Columns[x])
	}

	return &Board{
		w:           b.w,
		h:           b.h,
		winStreak:   b.winStreak,
		Columns:     columns,
		ColumnSizes: columnSizes,
	}
}

func ParseBoard(txt string, options ...Option) *Board {
	txt = strings.TrimSpace(txt)
	lines := strings.Split(txt, "\n")
	lines = ReverseLines(lines)
	for i, line := range lines {
		line = strings.ReplaceAll(line, " ", "")
		lines[i] = strings.ReplaceAll(line, "\t", "")
	}
	h := len(lines)
	w := len(lines[0])

	newOptions := append([]Option{WithSize(w, h)}, options...)
	board := NewBoard(newOptions...)
	for _, line := range lines {
		for x, cell := range line {
			if cell == rune(PlayerA) {
				board.Throw(x, PlayerA)
			} else if cell == rune(PlayerB) {
				board.Throw(x, PlayerB)
			}
		}
	}
	return board
}
