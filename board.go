package main

import (
	"fmt"
	"math/bits"
	"strings"

	"github.com/pkg/errors"
)

// Board keeps state in a binary format:
// - first binary One signifies stack height
// - next bits signifies Player A (zero) or Player B (one)
type Board struct {
	w         int
	h         int
	winStreak int
	state     BoardKey

	referee *Referee
}

type BoardKey [7]uint64

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
	b.state = [7]uint64{}
	for x := 0; x < b.w; x++ {
		b.state[x] = 0b1
	}
	b.referee = NewReferee(b)
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
func (b *Board) GetCell(x int, y int) Player {
	if b.state[x]>>y <= 1 {
		return Empty
	}
	return Player((b.state[x] >> y) & 0b1)
}

func (b *Board) Throw(x int, player Player) int {
	colSize := b.stackSize(x)
	// reset column signifying stack size
	b.state[x] = (b.state[x] & ^(1 << colSize)) | (1 << (colSize + 1)) | (uint64(player) << colSize)
	return colSize
}

func (b *Board) Revert(x int, y int) {
	// colsize = y + 1
	b.state[x] = (b.state[x] & ^(1 << (y + 1))) | (1 << y)
}

func (b *Board) stackSize(x int) int {
	return 7 - bits.LeadingZeros8(uint8(b.state[x]))
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
			if cell == PlayerA {
				rowCells = append(rowCells, ColouredPlayerA)
			} else if cell == PlayerB {
				rowCells = append(rowCells, ColouredPlayerB)
			} else {
				rowCells = append(rowCells, EmptyCell)
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

func (b *Board) HasWinner() Player {
	return b.referee.HasWinner(b)
}

func (b *Board) NextPlayer() Player {
	tokensA := 0
	tokensB := 0
	for x := 0; x < b.w; x++ {
		for y := 0; y < b.stackSize(x); y++ {
			cell := b.GetCell(x, y)
			if cell == PlayerA {
				tokensA++
			} else if cell == PlayerB {
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
	return b.state[x]>>b.h == 0
}

func (b *Board) Clone() *Board {
	state := [7]uint64{}
	for x := 0; x < b.w; x++ {
		state[x] = b.state[x]
	}

	return &Board{
		w:         b.w,
		h:         b.h,
		winStreak: b.winStreak,
		state:     state,
		referee:   b.referee,
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
			if cell == PlayerARune {
				board.Throw(x, PlayerA)
			} else if cell == PlayerBRune {
				board.Throw(x, PlayerB)
			}
		}
	}
	return board
}
