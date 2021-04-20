package common

import (
	"time"

	log "github.com/igrek51/log15"
)

const RefreshProgressPeriod = 2 * time.Second
const ItReportPeriodMask = 0b11111111111111111111 // modulo 2^20 (1048576) mask

type IMoveSolver interface {
	MovesEndings(board *Board) []Player
	HasPlayerWon(board *Board, move int, y int, player Player) bool
	SummaryVars() log.Ctx
	Cache() ICache
	Retrain(board *Board, maxDepth uint)
}
