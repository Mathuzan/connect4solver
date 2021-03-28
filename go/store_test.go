package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvaluateKey(t *testing.T) {
	board := ParseBoard(`
. . .
. B .
. A A
`)
	store := NewEndingStore()
	assert.EqualValues(t, "...AB.A..", store.EvaluateKey(board))
}
