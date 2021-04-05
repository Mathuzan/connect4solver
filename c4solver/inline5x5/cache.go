package inline5x5

import (
	"github.com/igrek51/connect4solver/c4solver/common"
	log "github.com/igrek51/log15"
)

const maxCacheSize = 1_500_000_000

type EndingCache struct {
	depthCaches       []map[uint64]common.Player
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
	depthCaches := make([]map[uint64]common.Player, boardW*boardH)
	for i := uint(0); i < uint(boardW*boardH); i++ {
		depthCaches[i] = make(map[uint64]common.Player)
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

func (s *EndingCache) Get(board *common.Board, depth uint) (ending common.Player, ok bool) {
	ending, ok = s.depthCaches[depth][s.reflectedBoardKey(board.State)]
	if !ok {
		return
	}
	s.cacheUsages++
	return
}

func (s *EndingCache) Put(board *common.Board, depth uint, ending common.Player) common.Player {
	if depth > s.maxCacheDepth {
		return ending
	}
	if len(s.depthCaches[depth]) >= s.maxCacheDepthSize {
		log.Debug("clearing cache", log.Ctx{"depth": depth})
		s.cachedEntries -= uint64(len(s.depthCaches[depth]))
		s.depthCaches[depth] = make(map[uint64]common.Player)
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

func (s *EndingCache) reflectedBoardKey(key common.BoardKey) uint64 {
	leftKey = key[0] | key[1]<<8
	rightKey = key[4] | key[3]<<8

	if leftKey <= rightKey {
		leftKey |= key[2] << 16
		leftKey |= key[3] << 24
		leftKey |= key[4] << 32
		return leftKey
	}
	// mirror map
	for i := 2; i < 5; i++ {
		rightKey |= key[2] << 16
		rightKey |= key[1] << 24
		rightKey |= key[0] << 32
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
