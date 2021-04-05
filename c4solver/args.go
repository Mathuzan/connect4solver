package c4solver

import (
	"flag"
	"fmt"

	"github.com/igrek51/connect4solver/c4solver/common"
)

func GetArgs() (int, int, int, bool, bool, common.Mode) {
	boardWidth := flag.Int("width", 4, "board width")
	boardHeight := flag.Int("height", 4, "board height")
	winStreak := flag.Int("win", 4, "win streak")
	boardSize := flag.String("size", "", "board size (7x6)")

	profileEnabled := flag.Bool("profile", false, "Enable pprof CPU profiling")
	cacheEnabled := flag.Bool("nocache", false, "Load cached endings from file")

	train := flag.Bool("train", false, "Training mode")
	play := flag.Bool("play", false, "Playing mode")

	flag.Parse()

	if boardSize != nil && *boardSize != "" {
		fmt.Sscanf(*boardSize, "%dx%d", boardWidth, boardHeight)
	}

	mode := common.TrainMode
	if *train {
		mode = common.TrainMode
	}
	if *play {
		mode = common.PlayMode
	}

	return *boardWidth, *boardHeight, *winStreak, *profileEnabled, !*cacheEnabled, mode
}
