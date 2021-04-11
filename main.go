package main

import (
	"os"
	"runtime/pprof"

	log "github.com/igrek51/log15"

	c4 "github.com/igrek51/connect4solver/solver"
	"github.com/igrek51/connect4solver/solver/common"
)

func main() {
	args := c4.GetArgs()

	if args.Profile {
		log.Info("Starting CPU profiler")
		cpuProfile, _ := os.Create("cpuprof.prof")
		pprof.StartCPUProfile(cpuProfile)
		defer pprof.StopCPUProfile()
	}

	if args.Mode == common.TrainMode {
		c4.Train(args.Width, args.Height, args.WinStreak, args.Cache)
	} else if args.Mode == common.PlayMode {
		c4.Play(args.Width, args.Height, args.WinStreak, args.Cache, args.HideA, args.HideB,
			args.AutoAttackA, args.AutoAttackB, args.Scores)
	}
}
