package c4solver

import (
	"fmt"
	"strings"
	"time"

	log "github.com/igrek51/log15"

	"github.com/igrek51/connect4solver/c4solver/common"
)

func Play(width, height, winStreak int, cacheEnabled bool, hideA, hideB bool) {
	board := common.NewBoard(common.WithSize(width, height), common.WithWinStreak(winStreak))

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

		if isATie(endings) {
			log.Info(fmt.Sprintf("%v", common.Tie))
			break
		}

		player := board.NextPlayer()
		fmt.Println(board.String())
		if (player == common.PlayerA && !hideA) || (player == common.PlayerB && !hideB) {
			printEndingsLine(endings, player)
		}

		move := readNextMove(endings, player)

		moveY := board.Throw(move, player)
		if solver.HasPlayerWon(board, move, moveY, player) {
			log.Info(fmt.Sprintf("Player %v won", player))
			break
		}
	}
}

func isATie(endings []common.Player) bool {
	for _, e := range endings {
		if e != common.NoMove {
			return false
		}
	}
	return true
}

func printEndingsLine(endings []common.Player, player common.Player) {
	displays := []string{}
	for _, ending := range endings {
		var display string
		if ending == common.NoMove {
			display = common.PlayerDisplays[common.NoMove]
		} else {
			playerEnding := EndingForPlayer(ending, player)
			display = common.ShortGameEndingDisplays[playerEnding]
		}
		displays = append(displays, display)
	}
	fmt.Println("| " + strings.Join(displays, " ") + " |")
}

func readNextMove(endings []common.Player, player common.Player) int {
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
		if endings[move] == common.NoMove {
			log.Error("Column is already full")
			continue
		}
		return move
	}
}
