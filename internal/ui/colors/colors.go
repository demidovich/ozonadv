package colors

import "github.com/fatih/color"

func Gray(s string, a ...any) string {
	return color.RGB(70, 70, 70).Sprintf(s, a...)
}

func Warning(s string, a ...any) string {
	return color.RGB(255, 255, 255).AddBgRGB(210, 0, 0).Sprintf(s, a...)
}
