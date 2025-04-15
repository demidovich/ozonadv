package helpers

import "github.com/fatih/color"

func TextGray(a ...any) string {
	return color.RGB(70, 70, 70).Sprint(a...)
}
