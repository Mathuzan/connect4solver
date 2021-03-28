package c4solver

import (
	"fmt"
	"time"

	log "github.com/inconshreveable/log15"
)

const BoardWidth = 7
const BoardHeight = 6

const MinWinCondition = 4

const PlayerA = 'A'
const PlayerB = 'B'

const CellA = PlayerA
const CellB = PlayerB
const CellEmpty = '.'

const Win = 'W'
const Lose = 'L'
const Tie = 'T'

var move_results_weights = map[rune]int{
	Win:  1,
	Tie:  0,
	Lose: -1,
}

func main() {
	width, height := getArgs()
	board := NewBoard(width, height)
	board.print()

	fmt.Println("Finding moves results...")
	startTime := time.Now()
	solver := MoveSolver()
	results := solver.MovesResults(board)
	for move, result := range results {
		log.Info("Move result", log.Ctx{"move": move, "result": result})
	}

	totalElapsed := time.Since(startTime)
	log.Info("Done", log.Ctx{"totalTime": totalElapsed})
}
