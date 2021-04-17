package common

import (
	log "github.com/igrek51/log15"
)

type IMoveSolver interface {
	MovesEndings(board *Board) []Player
	HasPlayerWon(board *Board, move int, y int, player Player) bool
	Interrupt()
	PreloadCache(board *Board)
	SaveCache()
	SummaryVars() log.Ctx
	Cache() ICache
	Retrain(board *Board, maxDepth uint)
}
