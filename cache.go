package main

type EndingCache struct {
	cache map[BoardKey]GameEnding
}

type BoardKey [7]byte

const OneByte byte = 1

func NewEndingCache() *EndingCache {
	return &EndingCache{
		cache: make(map[BoardKey]GameEnding),
	}
}

func (s *EndingCache) EvaluateKey(board *Board) BoardKey {
	var buffer [7]byte
	for x := 0; x < board.w; x++ {
		// first binary One signifies stack height
		buffer[x] = OneByte << board.ColumnSizes[x]
		for y := 0; y < board.ColumnSizes[x]; y++ {
			// zero means Player A, so already taken care by default
			if *board.Columns[x][y] == PlayerB {
				buffer[x] |= 1 << y
			}
		}
	}
	return buffer
}

func (s *EndingCache) Has(key BoardKey) bool {
	_, ok := s.cache[key]
	return ok
}

func (s *EndingCache) Get(key BoardKey) GameEnding {
	return s.cache[key]
}

func (s *EndingCache) Put(key BoardKey, ending GameEnding) {
	s.cache[key] = ending
}

func (s *EndingCache) Size() int {
	return len(s.cache)
}
