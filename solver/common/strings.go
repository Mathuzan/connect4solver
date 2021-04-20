package common

import (
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func ReverseLines(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

var numberPrinter *message.Printer = message.NewPrinter(language.Polish)

func BigintSeparated(num uint64) string {
	return strings.ReplaceAll(numberPrinter.Sprintf("%d", num), " ", "_")
}

func MaximumZeroIndex(nums []uint64) int {
	maxi := -1
	for i, num := range nums {
		if num == 0 {
			maxi = i
		} else {
			break
		}
	}
	return maxi
}
