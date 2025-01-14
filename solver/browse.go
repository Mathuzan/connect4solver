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
	startWithMoves string,
	retrainDepth int,
) {
	board := common.NewBoard(common.WithSize(width, height), common.WithWinStreak(winStreak))
	board.ApplyMoves(startWithMoves)

	solver := CreateSolver(board)
	if cacheEnabled && common.CacheFileExists(board) {
		common.MustLoadCache(solver.Cache(), board.W, board.H)
	}

	if retrainDepth > 0 {
		retrainSolverDepth(board, solver, uint(retrainDepth))
		common.MustSaveCache(solver.Cache(), board.W, board.H)
		return
	}

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
			solver.Cache().ClearCache(uint(x))
		} else if action == "clear_cache_from" {
			for d := uint(x); d < uint(board.W*board.H); d++ {
				solver.Cache().ClearCache(d)
			}
		} else if action == "endings" {
			startTime := time.Now()
			endings := solver.MovesEndings(board)
			if endings != nil {
				totalElapsed := time.Since(startTime)
				logger := log.New(log.Ctx{
					"solveTime": totalElapsed,
					"endings":   endings,
				})
				logger.Info("Board solved", solver.SummaryVars())
				printEndingsLine(endings, player)
			}
		} else if action == "cache" {
			showCacheStatistics(solver.Cache(), board.W, board.H)
			depth := board.CountMoves()
			cachedEndings := getCachedEndings(board, solver)
			log.Debug("cache statistics", log.Ctx{
				"depth":          depth,
				"depthCacheSize": solver.Cache().DepthSize(depth),
				"cachedEndings":  cachedEndings,
			})
			printGameEndingsLine(cachedEndings)
		} else if action == "save" {
			common.MustSaveCache(solver.Cache(), board.W, board.H)
		} else if action == "retrain" {
			retrainSolverDepth(board, solver, uint(x))
		}
	}
}

func readNextAction() (string, int) {
	for {
		fmt.Printf("Enter command (h for help) > ")
		var command string
		var x int

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
			fmt.Println("  clear X+ - clear cache from given depth")
			fmt.Println("  retrain X - retrain worst scenarios until given depth")
			fmt.Println("  save - save cache file")
			fmt.Println("  q - quit")
		} else if command == "" {
			return "", 0
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
			if command == "clear" {
				return "clear_cache_from", 0
			}
			if strings.HasSuffix(command, "+") {
				_, err := fmt.Sscanf(command, "clear %d+", &x)
				if err != nil {
					log.Error("Invalid number", log.Ctx{"error": err})
					continue
				}
				return "clear_cache_from", x
			} else {
				_, err := fmt.Sscanf(command, "clear %d", &x)
				if err != nil {
					log.Error("Invalid number", log.Ctx{"error": err})
					continue
				}
				return "clear_cache", x
			}
		} else if strings.HasPrefix(command, "retrain") {
			_, err := fmt.Sscanf(command, "retrain %d", &x)
			if err != nil {
				log.Error("Invalid number", log.Ctx{"error": err})
				continue
			}
			return "retrain", x
		} else if strings.HasPrefix(command, "m") {
			_, err := fmt.Sscanf(command, "m%d", &x)
			if err != nil {
				_, err2 := fmt.Sscanf(command, "m %d", &x)
				if err2 == nil {
					return "move", x
				}
				log.Error("Invalid number", log.Ctx{"error": err})
				continue
			}
			return "move", x
		} else if strings.HasPrefix(command, "r") {
			_, err := fmt.Sscanf(command, "r%d", &x)
			if err != nil {
				log.Error("Invalid number", log.Ctx{"error": err})
				continue
			}
			return "revert", x
		} else if move, err := strconv.Atoi(command); len(command) == 1 && err == nil {
			return "move", move
		} else {
			log.Error("Unknown command", log.Ctx{"command": command})
			continue
		}
	}
}

func getCachedEndings(board *common.Board, solver common.IMoveSolver) []common.GameEnding {
	endings := make([]common.GameEnding, board.W)
	player := board.NextPlayer()
	depth := board.CountMoves()
	for move := 0; move < board.W; move++ {
		if !board.CanMakeMove(move) {
			endings[move] = common.NoEnding
			continue
		}

		moveY := board.Throw(move, player)

		ending, ok := solver.Cache().Get(board, depth)
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

func retrainSolverDepth(board *common.Board, solver common.IMoveSolver, maxDepth uint) {
	log.Info("Retraining worst scenarios", log.Ctx{
		"maxDepth": maxDepth,
	})
	startTime := time.Now()
	solver.Retrain(board, maxDepth)
	totalElapsed := time.Since(startTime)
	logger := log.New(log.Ctx{
		"solveTime": totalElapsed,
	})
	logger.Info("Training done", solver.SummaryVars())
}

func showCacheStatistics(cache common.ICache, boardW, boardH int) {
	for d := uint(0); d < uint(boardW*boardH); d++ {
		log.Debug("depth cache", log.Ctx{
			"depth": d,
			"size":  cache.DepthSize(d),
		})
	}
}
