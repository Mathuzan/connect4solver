package main

type Player uint8

const (
	PlayerA Player = 0
	PlayerB Player = 1
	Empty   Player = 2

	PlayerARune = 'A'
	PlayerBRune = 'B'
	EmptyCell   = "."
)

type GameEnding string

const (
	Win  GameEnding = "Win"
	Tie  GameEnding = "Tie"
	Lose GameEnding = "Lose"
)

var PlayerDisplays = map[Player]string{
	PlayerA: "A",
	PlayerB: "B",
	Empty:   ".",
}

func (p Player) String() string {
	return PlayerDisplays[p]
}

func (e GameEnding) String() string {
	return string(e)
}
