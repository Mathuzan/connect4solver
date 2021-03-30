package main

import (
	"flag"
)

func getArgs() (int, int, int, bool) {
	boardWidth := flag.Int("width", 4, "board width")
	boardHeight := flag.Int("height", 4, "board height")
	winStreak := flag.Int("win", 4, "win streak")

	profileEnabled := flag.Bool("profile", false, "Enable pprof CPU profiling")

	flag.Parse()

	return *boardWidth, *boardHeight, *winStreak, *profileEnabled
}
