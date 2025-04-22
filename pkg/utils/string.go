package utils

import (
	"fmt"
	"math"
	"unicode/utf8"
)

func StringMasked(str string, length int) string {
	mask := "***"
	maskLen := 3
	minLen := maskLen + 2
	strLen := utf8.RuneCountInString(str)

	if strLen < minLen || length < minLen {
		return mask
	}

	segmentLen := int(math.Round(float64((length - maskLen) / 2)))
	r := []rune(str)

	return fmt.Sprintf(
		"%s%s%s",
		mask,
		string(r[:segmentLen]),
		string(r[strLen-segmentLen:]),
	)
}
