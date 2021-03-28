package main

import (
	"flag"
)

func getArgs() (int, int) {
	boardWidth := flag.Int("width", 4, "board width")
	boardHeight := flag.Int("height", 4, "board height")

	flag.Parse()

	return *boardWidth, *boardHeight
}
