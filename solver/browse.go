package solver

import (
	"fmt"
	"strings"
	"time"

	log "github.com/igrek51/log15"

	"github.com/igrek51/connect4solver/solver/common"
)

func Browse(
	width, height, winStreak int,
) {
	board := common.NewBoard(common.WithSize(width, height), common.WithWinStreak(winStreak))
	var solver *MoveSolver = NewMoveSolver(board)
	if CacheFileExists(board) {
		solver.PreloadCache(board)
	}

	for {
		fmt.Println(board.String())
		player := board.NextPlayer()
		fmt.Printf("Current player: %v\n", player)

		action, x := readNextAction()
		if action == "move" {
			board.Throw(x, player)
		} else if action == "revert" {
			board.Revert(x, board.StackSize(x)-1)
		} else if action == "endings" {
			startTime := time.Now()
			endings := solver.MovesEndings(board)
			totalElapsed := time.Since(startTime)
			logger := log.New(log.Ctx{
				"solveTime": totalElapsed,
			})
			logger.Info("Board solved", solver.SummaryVars())
			printEndingsLine(endings, player)
		} else if action == "cache" {
			cachedEndings := getCachedEndings(board, solver)
			printGameEndingsLine(cachedEndings)
			depth := board.CountMoves()
			depthCache := solver.cache.depthCaches[depth]
			log.Debug("cache statistics", log.Ctx{
				"depth":          depth,
				"depthCacheSize": len(depthCache),
			})
		}
	}
}

func readNextAction() (string, int) {
	for {
		fmt.Printf("Enter command (h for help) > ")
		var command string
		var move int
		fmt.Scanf("%s", &command)

		if command == "h" || command == "help" {
			fmt.Println("Available commands:")
			fmt.Println("  mX - move next player at column X, eg. m0")
			fmt.Println("  rX - revert token at column X, eg. r0")
			fmt.Println("  e - evaluate endings")
			fmt.Println("  c - show cached endings for current board")
		} else if strings.HasPrefix(command, "m") {
			_, err := fmt.Sscanf(command, "m%d", &move)
			if err != nil {
				log.Error("Invalid number", log.Ctx{"error": err})
				continue
			}
			return "move", move
		} else if strings.HasPrefix(command, "r") {
			_, err := fmt.Sscanf(command, "r%d", &move)
			if err != nil {
				log.Error("Invalid number", log.Ctx{"error": err})
				continue
			}
			return "revert", move
		} else if command == "e" {
			return "endings", 0
		} else if command == "c" {
			return "cache", 0
		} else {
			log.Error("Unknown command", log.Ctx{"command": command})
			continue
		}
	}
}

func getCachedEndings(board *common.Board, solver *MoveSolver) []common.GameEnding {
	endings := make([]common.GameEnding, board.W)
	player := board.NextPlayer()
	depth := board.CountMoves()
	for move := 0; move < board.W; move++ {
		if !board.CanMakeMove(move) {
			endings[move] = common.NoEnding
			continue
		}

		moveY := board.Throw(move, player)

		ending, ok := solver.cache.Get(board, depth)
		if !ok {
			endings[move] = common.NoEnding
		} else {
			playerEnding := common.EndingForPlayer(ending, player)
			endings[move] = playerEnding
		}

		board.Revert(move, moveY)
	}
	return endings
}

func printGameEndingsLine(endings []common.GameEnding) {
	displays := []string{}
	for _, ending := range endings {
		display := common.ShortGameEndingDisplays[ending]
		displays = append(displays, display)
	}
	fmt.Println("| " + strings.Join(displays, " ") + " |")
}
