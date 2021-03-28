package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertEqualTrimmed(t *testing.T, actual, expected string) {
	expectedTrimmed := strings.TrimSpace(expected)
	actual = strings.ReplaceAll(actual, "\u001b[31;1m", "")
	actual = strings.ReplaceAll(actual, "\u001b[33;1m", "")
	actual = strings.ReplaceAll(actual, "\u001b[0m", "")
	assert.Equal(t, expectedTrimmed, actual)
}
