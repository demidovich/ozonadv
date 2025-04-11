package prompts

import (
	"github.com/charmbracelet/huh"
)

func List(title string, options map[string]string) (string, error) {
	var value string

	o := make([]huh.Option[string], 0, len(options))
	for k, v := range options {
		o = append(o, huh.NewOption(k, v))
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(title).
				Options(o...).
				Value(&value),
		),
	)

	err := form.Run()

	return value, err
}
