package main

const BoardWidth = 7
const BoardHeight = 6

const MinWinCondition = 4

const CellEmpty = '.'

type Player rune

const (
	PlayerA Player = 'A'
	PlayerB Player = 'B'
)

type GameEnding rune

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
