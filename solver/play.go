package solver

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	log "github.com/igrek51/log15"

	"github.com/igrek51/connect4solver/solver/common"
)

func Play(
	width, height, winStreak int,
	cacheEnabled, hideA, hideB,
	autoAttackA, autoAttackB,
	scoresEnabled bool,
	startWithMoves string,
) {
	rand.Seed(time.Now().UnixNano())

	board := common.NewBoard(common.WithSize(width, height), common.WithWinStreak(winStreak))
	board.ApplyMoves(startWithMoves)

	solver := CreateSolver(board)
	if cacheEnabled && common.CacheFileExists(board) {
		common.MustLoadCache(solver.Cache(), board.W, board.H)
	}

	for {
		startTime := time.Now()
		endings := solver.MovesEndings(board)
		if endings == nil {
			endings = getCachedPlayerEndgames(board, solver)
		}

		player := board.NextPlayer()
		scores := estimateMoveScores(solver, endings, player, board, scoresEnabled)
		totalElapsed := time.Since(startTime)

		logger := log.New(log.Ctx{
			"solveTime": totalElapsed,
			"depth":     board.CountMoves(),
		})
		logger.Info("Board solved", solver.SummaryVars())

		fmt.Println(board.String())
		showHints := (player == common.PlayerA && !hideA) || (player == common.PlayerB && !hideB)
		if showHints {
			printEndingsLine(endings, player)
			if scoresEnabled {
				log.Info("Estimated move scores", log.Ctx{"scores": scores})
			}
		}
		bestMove := findBestMove(scores)

		var move int
		if (player == common.PlayerA && autoAttackA) || (player == common.PlayerB && autoAttackB) {
			move = bestMove
			playerEnding := common.EndingForPlayer(endings[move], player)
			fmt.Printf("Player %v moves: %d (%v)\n", player, move, playerEnding)
		} else {
			move = readNextMove(endings, player, board, bestMove, showHints)
		}

		moveY := board.Throw(move, player)
		if solver.HasPlayerWon(board, move, moveY, player) {
			depth := board.CountMoves()
			fmt.Println(board.String())
			log.Info(fmt.Sprintf("Player %v won in %d moves", player, depth))
			break
		} else if isATie(board) {
			depth := board.CountMoves()
			fmt.Println(board.String())
			log.Info(fmt.Sprintf("%v in %d moves", common.Tie, depth))
			break
		}
	}
}

func isATie(board *common.Board) bool {
	for x := 0; x < board.W; x++ {
		if board.CanMakeMove(x) {
			return false
		}
	}
	return true
}

func printEndingsLine(endings []common.Player, player common.Player) {
	if endings == nil {
		return
	}

	displays := []string{}
	for _, ending := range endings {
		var display string
		if ending == common.NoMove {
			display = common.PlayerDisplays[common.NoMove]
		} else {
			playerEnding := common.EndingForPlayer(ending, player)
			display = common.ShortGameEndingDisplays[playerEnding]
		}
		displays = append(displays, display)
	}
	fmt.Println("| " + strings.Join(displays, " ") + " |")
}

func readNextMove(
	endings []common.Player, player common.Player, board *common.Board,
	bestMove int, showBest bool,
) int {
	for {
		var move int
		bestStr := ""
		if showBest {
			bestStr = fmt.Sprintf(" (Best: %d)", bestMove)
		}
		fmt.Printf("Player %v moves [0-%d]%s: ", player, board.W-1, bestStr)
		_, err := fmt.Scanf("%d", &move)
		if err != nil {
			log.Error("Invalid number", log.Ctx{"error": err})
			continue
		}
		if move < 0 || move >= board.W {
			log.Error("Move number is out of range")
			continue
		}
		if endings != nil && !board.CanMakeMove(move) {
			log.Error("Column is already full")
			continue
		}
		return move
	}
}

func estimateMoveScores(
	solver common.IMoveSolver, endings []common.Player,
	player common.Player, board *common.Board, scoresEnabled bool,
) []int {
	scores := make([]int, len(endings))
	opponent := common.OppositePlayer(player)

	for move, ending := range endings {
		if ending == common.NoMove {
			scores[move] = -1000
			continue
		}

		moveY := board.Throw(move, player)
		score := 0

		if solver.HasPlayerWon(board, move, moveY, player) {
			score = 100
		} else if solver.HasPlayerWon(board, move, moveY, opponent) {
			score = -100
		} else {
			if ending == player {
				score += 10
			} else if ending == opponent {
				score -= 10
			}
			if scoresEnabled {
				nextEndings := solver.MovesEndings(board)
				for _, nextEnding := range nextEndings {
					if nextEnding == player {
						score++
					} else if nextEnding == opponent {
						score--
					}
				}
			}
		}

		scores[move] = score
		board.Revert(move, moveY)
	}
	return scores
}

func findBestMove(scores []int) int {
	order := rand.Perm(len(scores)) // get random if there are many maximum values
	maxi := order[0]
	for _, move := range order {
		if scores[move] > scores[maxi] {
			maxi = move
		}
	}
	return maxi
}

func getCachedPlayerEndgames(board *common.Board, solver common.IMoveSolver) []common.Player {
	endings := make([]common.Player, board.W)
	player := board.NextPlayer()
	depth := board.CountMoves()
	for move := 0; move < board.W; move++ {
		if !board.CanMakeMove(move) {
			endings[move] = common.NoMove
			continue
		}

		moveY := board.Throw(move, player)

		ending, ok := solver.Cache().Get(board, depth)
		if !ok {
			endings[move] = common.NoMove
		} else {
			endings[move] = ending
		}

		board.Revert(move, moveY)
	}
	return endings
}
