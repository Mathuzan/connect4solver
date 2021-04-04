package main

import (
	"flag"
	"fmt"
)

func getArgs() (int, int, int, bool, bool, Mode) {
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

	mode := TrainMode
	if *train {
		mode = TrainMode
	}
	if *play {
		mode = PlayMode
	}

	return *boardWidth, *boardHeight, *winStreak, *profileEnabled, !*cacheEnabled, mode
}
