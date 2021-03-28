package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertEqualTrimmed(t *testing.T, actual, expected string) {
	expectedTrimmed := strings.TrimSpace(expected)
	assert.Equal(t, expectedTrimmed, actual)
}
