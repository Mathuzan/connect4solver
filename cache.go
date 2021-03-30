package main

type EndingCache struct {
	cache map[BoardKey]GameEnding
}

func NewEndingCache() *EndingCache {
	return &EndingCache{
		cache: make(map[BoardKey]GameEnding),
	}
}

func (s *EndingCache) Has(board *Board) bool {
	_, ok := s.cache[board.state]
	return ok
}

func (s *EndingCache) Get(board *Board) GameEnding {
	return s.cache[board.state]
}

func (s *EndingCache) Put(board *Board, ending GameEnding) {
	s.cache[board.state] = ending
}

func (s *EndingCache) Size() int {
	return len(s.cache)
}
