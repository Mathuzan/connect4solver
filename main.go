package main

import (
	"os"
	"runtime/pprof"

	log "github.com/igrek51/log15"
)

var myPlayer = PlayerA

func main() {
	width, height, winStreak, profileEnabled, cacheEnabled, mode := getArgs()

	if profileEnabled {
		log.Info("Starting CPU profiler")
		cpuProfile, _ := os.Create("cpuprof.prof")
		pprof.StartCPUProfile(cpuProfile)
		defer pprof.StopCPUProfile()
	}

	if mode == TrainMode {
		Train(width, height, winStreak, cacheEnabled)
	} else if mode == PlayMode {
		Play(width, height, winStreak, cacheEnabled)
	}
}
