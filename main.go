package main

import (
	"os"
	"runtime/pprof"

	log "github.com/igrek51/log15"

	c4 "github.com/igrek51/connect4solver/c4solver"
	"github.com/igrek51/connect4solver/c4solver/common"
)

func main() {
	width, height, winStreak, profileEnabled, cacheEnabled, mode := c4.GetArgs()

	if profileEnabled {
		log.Info("Starting CPU profiler")
		cpuProfile, _ := os.Create("cpuprof.prof")
		pprof.StartCPUProfile(cpuProfile)
		defer pprof.StopCPUProfile()
	}

	if mode == common.TrainMode {
		c4.Train(width, height, winStreak, cacheEnabled)
	} else if mode == common.PlayMode {
		c4.Play(width, height, winStreak, cacheEnabled)
	}
}
