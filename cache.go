package main

type EndingCache struct {
	cache         map[BoardKey]GameEnding
	maxCacheDepth uint

	boardW  int
	boardW1 int
	sideW   int
}

func NewEndingCache(maxCacheDepth uint, boardW int) *EndingCache {
	return &EndingCache{
		cache:         make(map[BoardKey]GameEnding),
		maxCacheDepth: maxCacheDepth,
		boardW:        boardW,
		boardW1:       boardW - 1,
		sideW:         boardW / 2,
	}
}

func (s *EndingCache) Get(board *Board) (GameEnding, bool) {
	ending, ok := s.cache[s.reflectedBoardKey(board.state)]
	return ending, ok
}

func (s *EndingCache) Put(board *Board, depth uint, ending GameEnding) GameEnding {
	if depth > s.maxCacheDepth {
		return ending
	}
	s.cache[s.reflectedBoardKey(board.state)] = ending
	return ending
}

func (s *EndingCache) Size() int {
	return len(s.cache)
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
