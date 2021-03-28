package main

import (
	"flag"
)

func getArgs() (int, int, bool) {
	boardWidth := flag.Int("width", 4, "board width")
	boardHeight := flag.Int("height", 4, "board height")

	profileEnabled := flag.Bool("profile", false, "Enable pprof CPU profiling")

	flag.Parse()

	return *boardWidth, *boardHeight, *profileEnabled
}
