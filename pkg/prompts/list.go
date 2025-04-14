package prompts

import (
	"github.com/charmbracelet/huh"
)

type ListOption struct {
	Key   string
	Value string
}

func List(title string, options ...ListOption) (string, error) {
	var value string

	o := make([]huh.Option[string], 0, len(options))
	for _, v := range options {
		o = append(o, huh.NewOption(v.Key, v.Value))
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
