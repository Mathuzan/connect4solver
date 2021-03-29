package main

type EndingStore struct {
	cache map[string]GameEnding
}

func NewEndingStore() *EndingStore {
	return &EndingStore{
		cache: make(map[string]GameEnding),
	}
}

func (s *EndingStore) EvaluateKey(board *Board) string {
	buffer := make([]rune, board.w*board.h)
	var char rune
	for x := 0; x < board.w; x++ {
		for y := 0; y < board.h; y++ {
			if y >= board.ColumnSizes[x] {
				char = '.'
			} else {
				char = rune(*board.Columns[x][y])
			}
			buffer[y+x*board.h] = char
		}
	}
	return string(buffer)
}

func (s *EndingStore) Has(key string) bool {
	_, ok := s.cache[key]
	return ok
}

func (s *EndingStore) Get(key string) GameEnding {
	return s.cache[key]
}

func (s *EndingStore) Put(key string, ending GameEnding) {
	s.cache[key] = ending
}

func (s *EndingStore) Size() int {
	return len(s.cache)
}
