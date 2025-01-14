package common

import (
	"fmt"
	"math/bits"
	"strconv"
	"strings"

	log "github.com/igrek51/log15"
	"github.com/pkg/errors"
)

// Board keeps state in a binary format:
// - first binary One signifies stack height
// - next bits signifies Player A (zero) or Player B (one)
type Board struct {
	W         int
	H         int
	WinStreak int
	State     BoardKey
}

type BoardKey [7]uint64

// Board column state -> stack size
var StackSizeLookup [128]int

func init() {
	for state := uint8(1); state < 128; state++ {
		StackSizeLookup[state] = 7 - bits.LeadingZeros8(state)
	}
}

func NewBoard(options ...Option) *Board {
	// set defaults
	b := &Board{
		W:         4,
		H:         4,
		WinStreak: 4,
	}

	// apply options
	for _, opt := range options {
		err := opt(b)
		if err != nil {
			panic(errors.Wrap(err, "applying option"))
		}
	}

	b.State = [7]uint64{}
	b.Clear()
	return b
}

type Option func(*Board) error

func WithSize(width int, height int) Option {
	return func(b *Board) error {
		b.W = width
		b.H = height
		return nil
	}
}

func WithWinStreak(winStreak int) Option {
	return func(b *Board) error {
		b.WinStreak = winStreak
		return nil
	}
}

// GetCell return token at given cell coordinates (axes oriented top-right)
func (b *Board) GetCell(x int, y int) Player {
	if y >= b.StackSize(x) {
		return Empty
	}
	return Player((b.State[x] >> y) & 0b1)
}

func (b *Board) Throw(x int, player Player) int {
	colSize := b.StackSize(x)
	if player == PlayerA { // Player A = 0
		b.State[x] &= ^(1 << colSize) // reset column signifying old stack size
	}
	b.State[x] |= 1 << (colSize + 1)
	return colSize
}

func (b *Board) Revert(x int, y int) {
	// colsize = y + 1
	b.State[x] = (b.State[x] & ^(1 << (y + 1))) | (1 << y)
}

func (b *Board) StackSize(x int) int {
	return StackSizeLookup[b.State[x]]
}

func (b *Board) String() string {
	lines := []string{}
	var line string

	line = "+-"
	for i := 0; i < b.W; i++ {
		line += "--"
	}
	line += "+"
	lines = append(lines, line)

	for y := b.H - 1; y >= 0; y-- {
		rowCells := []string{}
		for x := 0; x < b.W; x++ {
			cell := b.GetCell(x, y)
			rowCells = append(rowCells, PlayerDisplays[cell])
		}
		line = "| " + strings.Join(rowCells, " ") + " |"
		lines = append(lines, line)
	}

	line = "+-"
	for i := 0; i < b.W; i++ {
		line += "--"
	}
	line += "+"
	lines = append(lines, line)

	coordinates := []string{}
	for x := 0; x < b.W; x++ {
		coordinates = append(coordinates, fmt.Sprint(x))
	}
	line = "| " + strings.Join(coordinates, " ") + " |"
	lines = append(lines, line)

	return strings.Join(lines, "\n")
}

func (b *Board) NextPlayer() Player {
	tokensA := 0
	tokensB := 0
	for x := 0; x < b.W; x++ {
		for y := 0; y < b.StackSize(x); y++ {
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

func (b *Board) CountMoves() uint {
	tokens := uint(0)
	for x := 0; x < b.W; x++ {
		for y := 0; y < b.StackSize(x); y++ {
			cell := b.GetCell(x, y)
			if cell == PlayerA || cell == PlayerB {
				tokens++
			}
		}
	}
	return tokens
}

func (b *Board) CanMakeMove(x int) bool {
	return b.State[x]>>b.H == 0
}

func (b *Board) Clear() {
	for x := 0; x < b.W; x++ {
		b.State[x] = 0b1
	}
}

func (b *Board) Clone() *Board {
	state := [7]uint64{}
	for x := 0; x < b.W; x++ {
		state[x] = b.State[x]
	}

	return &Board{
		W:         b.W,
		H:         b.H,
		WinStreak: b.WinStreak,
		State:     state,
	}
}

func (b *Board) ApplyMoves(startWithMoves string) *Board {
	if startWithMoves == "" {
		return b
	}
	for idx, moveStr := range startWithMoves {
		move, err := strconv.Atoi(string(moveStr))
		if err != nil {
			log.Error("Invalid number", log.Ctx{"error": err, "index": idx, "move": moveStr})
			return b
		}
		if move < 0 || move >= b.W {
			log.Error("Move number is out of range", log.Ctx{"move": move, "index": idx})
			return b
		}
		if !b.CanMakeMove(move) {
			log.Error("Column is already full", log.Ctx{"move": move, "index": idx})
			return b
		}
		player := b.NextPlayer()
		b.Throw(move, player)
	}
	return b
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
