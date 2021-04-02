package main

type EndingCache struct {
	depthCaches   []map[uint64]GameEnding
	maxCacheDepth uint
	cachedEntries uint64

	boardW  int
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
		boardW1:       boardW - 1,
		sideW:         boardW / 2,
	}
}

func (s *EndingCache) Get(board *Board, depth uint) (ending GameEnding, ok bool) {
	ending, ok = s.depthCaches[depth][s.reflectedBoardKey(board.state)]
	return
}

func (s *EndingCache) Put(board *Board, depth uint, ending GameEnding) GameEnding {
	if depth > s.maxCacheDepth {
		return ending
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
