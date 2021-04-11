package solver

import (
	"flag"
	"fmt"

	"github.com/igrek51/connect4solver/solver/common"
)

type CliArgs struct {
	Width     int
	Height    int
	WinStreak int

	Mode common.Mode

	Profile     bool
	Cache       bool
	HideA       bool
	HideB       bool
	AutoAttackA bool
	AutoAttackB bool
	Scores      bool
}

func GetArgs() *CliArgs {
	args := &CliArgs{}

	flag.IntVar(&args.Width, "width", 7, "board width")
	flag.IntVar(&args.Height, "height", 6, "board height")
	flag.IntVar(&args.WinStreak, "win", 4, "win streak")
	boardSize := flag.String("size", "", "board size (eg. 7x6)")

	flag.BoolVar(&args.Profile, "profile", false, "Enable pprof CPU profiling")
	nocache := flag.Bool("nocache", false, "Load cached endings from file")
	flag.BoolVar(&args.HideA, "hide-a", false, "Hide endings hints for player A")
	flag.BoolVar(&args.HideB, "hide-b", false, "Hide endings hints for player B")
	flag.BoolVar(&args.AutoAttackA, "autoattack-a", false, "Make player A move automatically")
	flag.BoolVar(&args.AutoAttackB, "autoattack-b", false, "Make player B move automatically")
	flag.BoolVar(&args.Scores, "scores", false, "Show scores of each move")

	train := flag.Bool("train", false, "Training mode")
	play := flag.Bool("play", false, "Playing mode")

	flag.Parse()

	if boardSize != nil && *boardSize != "" {
		fmt.Sscanf(*boardSize, "%dx%d", &args.Width, &args.Height)
	}

	args.Cache = !*nocache

	args.Mode = common.TrainMode
	if *train {
		args.Mode = common.TrainMode
	}
	if *play {
		args.Mode = common.PlayMode
	}

	return args
}