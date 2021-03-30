package main

type EndingCache struct {
	depthCaches   map[uint]map[BoardKey]GameEnding
	maxCacheDepth uint
	cachedEntries uint64

	boardW  int
	boardW1 int
	sideW   int
}

func NewEndingCache(maxCacheDepth uint, boardW int, boardH int) *EndingCache {
	depthCaches := make(map[uint]map[BoardKey]GameEnding)
	for i := uint(0); i < uint(boardW*boardH); i++ {
		depthCaches[i] = make(map[BoardKey]GameEnding)
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

func (s *EndingCache) reflectedBoardKey(key BoardKey) BoardKey {
	var leftKey uint = 0
	var rightKey uint = 0
	for i := 0; i < s.sideW; i++ {
		leftKey += uint(key[i]) << (8 * i)
		rightKey += uint(key[s.boardW1-i]) << (8 * i)
	}

	if leftKey <= rightKey {
		return key
	}

	var newKey BoardKey
	for i := 0; i < s.boardW; i++ {
		newKey[i] = key[s.boardW1-i]
	}

	return newKey
}
