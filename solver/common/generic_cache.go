package common

type ICache interface {
	Get(board *Board, depth uint) (ending Player, ok bool)
	Put(board *Board, depth uint, ending Player) Player
	ClearCache(depth uint)
	Size() uint64
	DepthSize(depth uint) int
	ShowStatistics()
}
