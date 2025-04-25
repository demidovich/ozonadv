package utils

import (
	"fmt"
	"math"
	"unicode/utf8"
)

func StringMasked(str, mask string, maxlen int) string {
	strlen := utf8.RuneCountInString(str)
	masklen := utf8.RuneCountInString(mask)

	if strlen < maxlen {
		maxlen = strlen
	}

	if masklen < 1 {
		return str
	}

	if maxlen < 1 || strlen < masklen+2 {
		return mask
	}

	var halflen int
	if maxlen >= masklen+2 {
		halflen = int(math.Round(float64((maxlen - masklen) / 2)))
	} else {
		halflen = int(math.Round(float64((masklen - 2) / 2)))
	}

	r := []rune(str)
	return fmt.Sprintf(
		"%s%s%s",
		string(r[:halflen]),
		mask,
		string(r[strlen-halflen:]),
	)
}
