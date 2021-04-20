package inline7x6

import (
	"fmt"
	"time"

	log "github.com/igrek51/log15"

	"github.com/igrek51/connect4solver/solver/common"
)

func (s *MoveSolver) reportCycle(board *common.Board, progressStart float64) {
	if s.iterations&common.ItReportPeriodMask == 0 && time.Since(s.lastBoardPrintTime) >= common.RefreshProgressPeriod {
		s.ReportStatus(board, progressStart)
		if s.interrupt {
			panic(common.InterruptError)
		}
	}
}

func (s *MoveSolver) ReportStatus(
	board *common.Board,
	progress float64,
) {
	duration := time.Since(s.startTime)
	instDuration := time.Since(s.lastBoardPrintTime)
	if s.firstProgress == 0 {
		s.firstProgress = progress
	}
	var eta time.Duration
	if progress > s.firstProgress && duration > 0 {
		eta = time.Duration((1-progress)*1000/(progress-s.firstProgress)) * (duration - common.RefreshProgressPeriod) / 1000
	}
	iterationsPerSec := common.BigintSeparated(s.iterations / uint64(duration/time.Second))
	instIterationsPerSec := common.BigintSeparated((s.iterations - s.lastIterations) / uint64(instDuration/time.Second))
	s.lastIterations = s.iterations
	s.lastBoardPrintTime = time.Now()

	firstCaches := []int{}
	for d := uint(0); d < uint(10); d++ {
		firstCaches = append(firstCaches, s.Cache().DepthSize(d))
	}

	log.Debug("Currently considered board", log.Ctx{
		"cacheSize":         common.BigintSeparated(s.cache.Size()),
		"iterations":        common.BigintSeparated(s.iterations),
		"cacheClears":       common.BigintSeparated(s.cache.clears),
		"maxUnclearedDepth": common.MaximumZeroIndex(s.cache.depthClears),
		"progress":          fmt.Sprintf("%v", progress),
		"eta":               eta.Truncate(time.Second),
		"itsAvg":            iterationsPerSec,
		"its":               instIterationsPerSec,
		"firstCacheLen":     firstCaches,
	})
	fmt.Println(board.String())
	if s.progressBar != nil {
		s.progressBar.Set(int(progress * common.ProgressBarResolution))
	}
}

func (s *MoveSolver) SummaryVars() log.Ctx {
	return log.Ctx{
		"cacheSize":   common.BigintSeparated(s.cache.Size()),
		"iterations":  common.BigintSeparated(s.iterations),
		"cacheClears": common.BigintSeparated(s.cache.clears),
	}
}
