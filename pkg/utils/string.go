package utils

import (
	"fmt"
	"math"
	"unicode/utf8"
)

func StringMasked(str string, length int) string {
	strLen := utf8.RuneCountInString(str)
	if strLen < 5 || length < 5 {
		return "***"
	}

	segmentLen := int(math.Round(float64((length - 3) / 2)))
	r := []rune(str)

	return fmt.Sprintf(
		"%s***%s",
		string(r[:segmentLen]),
		string(r[strLen-segmentLen:]),
	)
}
