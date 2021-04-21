package generic_solver

import (
	"github.com/igrek51/connect4solver/solver/common"
	log "github.com/igrek51/log15"
)

const maxCacheSize = 1_500_000_000

type EndingCache struct {
	depthCaches            []map[uint64]common.Player
	maxCacheDepthSize      int
	maxCachedDepth         uint
	maxUnclearedCacheDepth uint

	cachedEntries uint64
	cacheUsages   uint64
	clears        uint64
	depthClears   []uint64

	boardW  int
	boardH  int
	boardW1 int
	sideW   int
}

func NewEndingCache(boardW int, boardH int) *EndingCache {
	depthCaches := make([]map[uint64]common.Player, boardW*boardH)
	for i := uint(0); i < uint(boardW*boardH); i++ {
		depthCaches[i] = make(map[uint64]common.Player)
	}

	return &EndingCache{
		depthCaches:            depthCaches,
		depthClears:            make([]uint64, boardW*boardH),
		maxCacheDepthSize:      maxCacheSize / (boardW * boardH),
		maxCachedDepth:         uint(boardW*boardH) - 4,
		maxUnclearedCacheDepth: 16,
		boardW:                 boardW,
		boardH:                 boardH,
		boardW1:                boardW - 1,
		sideW:                  boardW / 2,
	}
}

func (s *EndingCache) Get(board *common.Board, depth uint) (ending common.Player, ok bool) {
	ending, ok = s.depthCaches[depth][s.reflectedBoardKey(board.State)]
	return
}

func (s *EndingCache) Put(board *common.Board, depth uint, ending common.Player) common.Player {
	if depth > s.maxCachedDepth {
		return ending
	}
	if s.DepthSize(depth) >= s.maxCacheDepthSize && depth > s.maxUnclearedCacheDepth {
		s.ClearCache(depth)
	}
	s.depthCaches[depth][s.reflectedBoardKey(board.State)] = ending
	s.cachedEntries++
	return ending
}

func (s *EndingCache) ClearCache(depth uint) {
	log.Debug("clearing cache", log.Ctx{"depth": depth})
	s.cachedEntries -= uint64(s.DepthSize(depth))
	s.depthCaches[depth] = make(map[uint64]common.Player)
	s.depthClears[depth]++
	s.clears++
}

func (s *EndingCache) Size() uint64 {
	return s.cachedEntries
}

func (s *EndingCache) DepthSize(depth uint) int {
	return len(s.depthCaches[depth])
}

func (s *EndingCache) reflectedBoardKey(key common.BoardKey) uint64 {
	leftKey := key[0]
	rightKey := key[s.boardW1]
	for i := 1; i < s.sideW; i++ {
		leftKey |= key[i] << (8 * i)
		rightKey |= key[s.boardW1-i] << (8 * i)
	}

	if leftKey <= rightKey {
		for i := s.sideW; i < s.boardW; i++ {
			leftKey |= key[i] << (8 * i)
		}
		return leftKey
	}
	// mirror map
	for i := s.sideW; i < s.boardW; i++ {
		rightKey |= key[s.boardW1-i] << (8 * i)
	}
	return rightKey
}

func (s *EndingCache) MaxCachedDepth() uint {
	return s.maxCachedDepth
}

func (s *EndingCache) DepthCaches() []map[uint64]common.Player {
	return s.depthCaches
}

func (s *EndingCache) SetEntry(depth int, key uint64, value common.Player) {
	s.depthCaches[depth][key] = value
	s.cachedEntries++
}
