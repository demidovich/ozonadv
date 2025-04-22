package colors

import "github.com/fatih/color"

func Gray() *color.Color {
	return color.RGB(70, 70, 70)
}

func Warning() *color.Color {
	return color.RGB(255, 255, 255).AddBgRGB(210, 0, 0)
}
