package main

type EndingStore struct {
	cache map[string]GameEnding
}

func NewEndingStore() *EndingStore {
	return &EndingStore{
		cache: make(map[string]GameEnding),
	}
}
