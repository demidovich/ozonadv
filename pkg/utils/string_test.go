package utils

import (
	"strconv"
	"testing"
)

func TestStringMasked(t *testing.T) {
	var tests = []struct {
		input  string
		mask   string
		maxlen int
		want   string
	}{
		{"", "***", 5, "***"},

		{"12345", "***", -1, "***"},
		{"12345", "***", 0, "***"},
		{"12345", "***", 1, "***"},
		{"12345", "***", 4, "***"},
		{"12345", "***", 5, "1***5"},
		{"12345", "***", 7, "1***5"},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := StringMasked(tt.input, tt.mask, tt.maxlen)
			if actual != tt.want {
				t.Errorf("got %s, want %s", actual, tt.want)
			}
		})
	}
}
