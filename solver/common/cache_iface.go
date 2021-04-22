package common

var CacheSizeLimit int = 1_500_000_000

type ICache interface {
	Get(board *Board, depth uint) (ending Player, ok bool)
	Put(board *Board, depth uint, ending Player) Player
	ClearCache(depth uint)
	Size() uint64
	DepthSize(depth uint) int
	MaxCachedDepth() uint
	DepthCaches() []map[uint64]Player
	SetEntry(depth int, key uint64, value Player)
}
