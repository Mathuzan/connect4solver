package c4solver

import (
	. "github.com/igrek51/connect4solver/c4solver/common"
	log "github.com/igrek51/log15"
)

const maxCacheSize = 1_500_000_000

type EndingCache struct {
	depthCaches       []map[uint64]Player
	maxCacheDepthSize int
	maxCacheDepth     uint
	cachedEntries     uint64
	cacheUsages       uint64
	clears            uint64

	boardW  int
	boardH  int
	boardW1 int
	sideW   int
}

func NewEndingCache(boardW int, boardH int) *EndingCache {
	depthCaches := make([]map[uint64]Player, boardW*boardH)
	for i := uint(0); i < uint(boardW*boardH); i++ {
		depthCaches[i] = make(map[uint64]Player)
	}

	return &EndingCache{
		depthCaches:       depthCaches,
		maxCacheDepthSize: maxCacheSize / (boardW * boardH),
		maxCacheDepth:     uint(boardW*boardH) - 4,
		boardW:            boardW,
		boardH:            boardH,
		boardW1:           boardW - 1,
		sideW:             boardW / 2,
	}
}

func (s *EndingCache) Get(board *Board, depth uint) (ending Player, ok bool) {
	ending, ok = s.depthCaches[depth][s.reflectedBoardKey(board.State)]
	if !ok {
		return
	}
	s.cacheUsages++
	return
}

func (s *EndingCache) Put(board *Board, depth uint, ending Player) Player {
	if depth > s.maxCacheDepth {
		return ending
	}
	if len(s.depthCaches[depth]) >= s.maxCacheDepthSize {
		log.Debug("clearing cache", log.Ctx{"depth": depth})
		s.cachedEntries -= uint64(len(s.depthCaches[depth]))
		s.depthCaches[depth] = make(map[uint64]Player)
		s.clears++
	}
	s.depthCaches[depth][s.reflectedBoardKey(board.State)] = ending
	s.cachedEntries++
	return ending
}

func (s *EndingCache) Size() uint64 {
	return s.cachedEntries
}

var leftKey uint64 = 0
var rightKey uint64 = 0

func (s *EndingCache) reflectedBoardKey(key BoardKey) uint64 {
	leftKey = key[0]
	rightKey = key[s.boardW1]
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

func (s *EndingCache) ShowStatistics() {
	for d := uint(0); d < uint(s.boardW*s.boardH); d++ {
		depthCache := s.depthCaches[d]

		log.Debug("depth cache", log.Ctx{
			"depth": d,
			"size":  len(depthCache),
		})
	}
}

func (s *EndingCache) HighestDepth() int {
	maxd := 0
	for d, depthCache := range s.depthCaches {
		if len(depthCache) > len(s.depthCaches[maxd]) {
			maxd = d
		}
	}
	return maxd
}