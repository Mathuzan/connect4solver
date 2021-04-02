package main

const BoardWidth = 7
const BoardHeight = 6

const MinWinCondition = 4

type Player uint8

const (
	PlayerA Player = 0
	PlayerB Player = 1
	Empty   Player = 2

	PlayerARune = 'A'
	PlayerBRune = 'B'
	EmptyCell   = "."
)

type GameEnding uint8

const (
	Win  GameEnding = 'W'
	Tie  GameEnding = 'T'
	Lose GameEnding = 'L'
)

var MoveResultsWeights = map[GameEnding]int{
	Win:  1,
	Tie:  0,
	Lose: -1,
}

var PlayerDisplays = map[Player]rune{
	PlayerA: 'A',
	PlayerB: 'B',
	Empty:   '.',
}

func (p Player) String() string {
	return string(p)
}

func (e GameEnding) String() string {
	return string(e)
}
