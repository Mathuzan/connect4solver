package main

import (
	"fmt"
	"strings"
	"time"

	log "github.com/igrek51/log15"
)

func Play(width, height, winStreak int, cacheEnabled bool) {
	board := NewBoard(WithSize(width, height), WithWinStreak(winStreak))

	var solver IMoveSolver = NewMoveSolver(board)

	if cacheEnabled && CacheFileExists(board) {
		solver.PreloadCache(board)
	}

	for {
		startTime := time.Now()
		endings := solver.MovesEndings(board)
		totalElapsed := time.Since(startTime)

		logger := log.New(log.Ctx{
			"solveTime":   totalElapsed,
			"boardWidth":  width,
			"boardHeight": height,
			"winStreak":   winStreak,
		})
		logger.Info("Board solved", solver.ContextVars())

		player := board.NextPlayer()
		fmt.Println(board.String())
		printEndingsLine(endings, player)

		move := readNextMove(endings, player)

		moveY := board.Throw(move, player)
		if board.referee.HasPlayerWon(board, move, moveY, player) {
			log.Info(fmt.Sprintf("Player %v won", player))
			break
		} else if isATie(endings) {
			log.Info(fmt.Sprintf("%v", Tie))
			break
		}
	}
}

func isATie(endings []Player) bool {
	for _, e := range endings {
		if e != NoMove {
			return false
		}
	}
	return true
}

func printEndingsLine(endings []Player, player Player) {
	displays := []string{}
	for _, ending := range endings {
		var display string
		if ending == NoMove {
			display = PlayerDisplays[NoMove]
		} else {
			playerEnding := EndingForPlayer(ending, player)
			display = ShortGameEndingDisplays[playerEnding]
		}
		displays = append(displays, display)
	}
	fmt.Println("| " + strings.Join(displays, " ") + " |")
}

func readNextMove(endings []Player, player Player) int {
	for {
		var move int
		fmt.Printf("Player %v moves [0-%d]: ", player, len(endings)-1)
		_, err := fmt.Scanf("%d", &move)
		if err != nil {
			log.Error("Invalid number", log.Ctx{"error": err})
			continue
		}
		if move < 0 || move >= len(endings) {
			log.Error("Move number is out of range")
			continue
		}
		if endings[move] == NoMove {
			log.Error("Cant move at full column")
			continue
		}
		return move
	}
}
