package solver

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/igrek51/log15"

	"github.com/igrek51/connect4solver/solver/common"
)

func Browse(
	width, height, winStreak int,
	cacheEnabled bool,
) {
	board := common.NewBoard(common.WithSize(width, height), common.WithWinStreak(winStreak))
	var solver *MoveSolver = NewMoveSolver(board)
	if cacheEnabled && CacheFileExists(board) {
		solver.PreloadCache(board)
	}
	common.HandleInterrupt(solver)

	for {
		fmt.Println(board.String())
		player := board.NextPlayer()
		depth := board.CountMoves()
		fmt.Printf("Current player: %v, moves: %d\n", player, depth)

		action, x := readNextAction()
		if action == "quit" {
			return
		} else if action == "move" {
			if x < 0 || x >= board.W {
				log.Error("Move number is out of range")
				continue
			}
			if !board.CanMakeMove(x) {
				log.Error("Column is already full")
				continue
			}
			board.Throw(x, player)
		} else if action == "revert" {
			if x < 0 || x >= board.W {
				log.Error("Move number is out of range")
				continue
			}
			if board.StackSize(x) == 0 {
				log.Error("Column is already empty")
				continue
			}
			board.Revert(x, board.StackSize(x)-1)
		} else if action == "new" {
			board.Clear()
		} else if action == "clear_cache" {
			solver.cache.ClearCache(uint(x))
		} else if action == "endings" {
			startTime := time.Now()
			if solver.interrupt {
				common.HandleInterrupt(solver)
			}
			endings := solver.MovesEndings(board)
			totalElapsed := time.Since(startTime)
			logger := log.New(log.Ctx{
				"solveTime": totalElapsed,
				"endings":   endings,
			})
			if endings != nil {
				logger.Info("Board solved", solver.SummaryVars())
				printEndingsLine(endings, player)
			}
		} else if action == "cache" {
			solver.cache.ShowStatistics()
			depth := board.CountMoves()
			cachedEndings := getCachedEndings(board, solver)
			depthCache := solver.cache.depthCaches[depth]
			log.Debug("cache statistics", log.Ctx{
				"depth":          depth,
				"depthCacheSize": len(depthCache),
				"cachedEndings":  cachedEndings,
			})
			printGameEndingsLine(cachedEndings)
		} else if action == "save" {
			solver.SaveCache()
		}
	}
}

func readNextAction() (string, int) {
	for {
		fmt.Printf("Enter command (h for help) > ")
		var command string
		var move int

		in := bufio.NewReader(os.Stdin)
		command, err := in.ReadString('\n')
		if err != nil {
			log.Error("Command read error", log.Ctx{"command": command})
			continue
		}
		command = strings.TrimSuffix(command, "\n")

		if command == "h" || command == "help" {
			fmt.Println("Available commands:")
			fmt.Println("  X, mX - move next player at column X [0-6], eg. m0")
			fmt.Println("  rX - revert token at column X, eg. r0")
			fmt.Println("  e - evaluate endings")
			fmt.Println("  c - show cache statistics & cached endings for current board")
			fmt.Println("  new - start new game")
			fmt.Println("  clear X - clear cache at given depth")
			fmt.Println("  save - save cache file")
			fmt.Println("  q - quit")
		} else if command == "" {
			return "", 0
		} else if strings.HasPrefix(command, "m") {
			_, err := fmt.Sscanf(command, "m%d", &move)
			if err != nil {
				_, err2 := fmt.Sscanf(command, "m %d", &move)
				if err2 == nil {
					return "move", move
				}
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
		} else if command == "q" {
			return "quit", 0
		} else if command == "e" {
			return "endings", 0
		} else if command == "c" {
			return "cache", 0
		} else if command == "new" {
			return "new", 0
		} else if command == "save" {
			return "save", 0
		} else if strings.HasPrefix(command, "clear") {
			_, err := fmt.Sscanf(command, "clear %d", &move)
			if err != nil {
				log.Error("Invalid number", log.Ctx{"error": err})
				continue
			}
			return "clear_cache", move
		} else if move, err := strconv.Atoi(command); len(command) == 1 && err == nil {
			return "move", move
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
