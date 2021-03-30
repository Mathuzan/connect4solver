package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvaluateKey(t *testing.T) {
	board := ParseBoard(`
. . . .
. . . B
. A B A
`)
	key := board.state
	byte1 := [7]byte(key)[0]
	byte2 := [7]byte(key)[1]
	byte3 := [7]byte(key)[2]
	byte4 := [7]byte(key)[3]
	assert.EqualValues(t, byte(0b00000001), byte1)
	assert.EqualValues(t, byte(0b00000010), byte2)
	assert.EqualValues(t, byte(0b00000011), byte3)
	assert.EqualValues(t, byte(0b00000110), byte4)
}

func TestMirrorKey(t *testing.T) {
	boardL := ParseBoard(`
A . . . .
A . . B .
A . B A A
`)
	boardR := ParseBoard(`
. . . . A
. B . . A
A A B . A
`)
	cache := NewEndingCache(999, 5)

	cache.Put(boardL, 7, Win)

	end, ok := cache.Get(boardL)
	assert.EqualValues(t, true, ok)
	assert.EqualValues(t, Win, end)

	end, ok = cache.Get(boardR)
	assert.EqualValues(t, true, ok)
	assert.EqualValues(t, Win, end)
}
