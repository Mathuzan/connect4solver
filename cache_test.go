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
	store := NewEndingCache()
	key := store.EvaluateKey(board)
	byte1 := [7]byte(key)[0]
	byte2 := [7]byte(key)[1]
	byte3 := [7]byte(key)[2]
	byte4 := [7]byte(key)[3]
	assert.EqualValues(t, byte(0b00000001), byte1)
	assert.EqualValues(t, byte(0b00000010), byte2)
	assert.EqualValues(t, byte(0b00000011), byte3)
	assert.EqualValues(t, byte(0b00000110), byte4)
}
