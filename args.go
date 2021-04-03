package main

import (
	"flag"
	"fmt"
)

func getArgs() (int, int, int, bool, bool) {
	boardWidth := flag.Int("width", 4, "board width")
	boardHeight := flag.Int("height", 4, "board height")
	winStreak := flag.Int("win", 4, "win streak")
	boardSize := flag.String("size", "", "board size (7x6)")

	profileEnabled := flag.Bool("profile", false, "Enable pprof CPU profiling")
	cacheEnabled := flag.Bool("cache", false, "Load cached endings from file")

	flag.Parse()

	if boardSize != nil && *boardSize != "" {
		fmt.Sscanf(*boardSize, "%dx%d", boardWidth, boardHeight)
	}

	return *boardWidth, *boardHeight, *winStreak, *profileEnabled, *cacheEnabled
}
