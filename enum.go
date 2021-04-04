package main

type Player uint8

const (
	PlayerA Player = 0 // Player A moves first
	PlayerB Player = 1 // Player B has a second move
	Empty   Player = 2

	PlayerARune = 'A'
	PlayerBRune = 'B'
	EmptyCell   = "."
)

type GameEnding string

const (
	Win  GameEnding = "\u001b[32;1mWin\u001b[0m"
	Tie  GameEnding = "\u001b[33;1mTie\u001b[0m"
	Lose GameEnding = "\u001b[31;1mLose\u001b[0m"
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
