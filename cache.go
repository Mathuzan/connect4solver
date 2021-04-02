package main

import (
	log "github.com/igrek51/log15"
)

const maxCacheDepthSize = 200_000_000

type EndingCache struct {
	depthCaches   []map[uint64]GameEnding
	maxCacheDepth uint
	cachedEntries uint64
	cacheUsages   uint64

	boardW  int
	boardH  int
	boardW1 int
	sideW   int
}

func NewEndingCache(maxCacheDepth uint, boardW int, boardH int) *EndingCache {
	depthCaches := make([]map[uint64]GameEnding, boardW*boardH)
	for i := uint(0); i < uint(boardW*boardH); i++ {
		depthCaches[i] = make(map[uint64]GameEnding)
	}

	return &EndingCache{
		depthCaches:   depthCaches,
		maxCacheDepth: maxCacheDepth,
		boardW:        boardW,
		boardH:        boardH,
		boardW1:       boardW - 1,
		sideW:         boardW / 2,
	}
}

func (s *EndingCache) Get(board *Board, depth uint) (ending GameEnding, ok bool) {
	ending, ok = s.depthCaches[depth][s.reflectedBoardKey(board.state)]
	if !ok {
		return
	}
	s.cacheUsages++
	return
}

func (s *EndingCache) Put(board *Board, depth uint, ending GameEnding) GameEnding {
	if depth > s.maxCacheDepth {
		return ending
	}
	if len(s.depthCaches[depth]) >= maxCacheDepthSize {
		log.Debug("clearing cache", log.Ctx{"depth": depth})
		s.depthCaches[depth] = make(map[uint64]GameEnding)
	}
	s.depthCaches[depth][s.reflectedBoardKey(board.state)] = ending
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
