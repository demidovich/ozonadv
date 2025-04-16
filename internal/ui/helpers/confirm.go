package helpers

import "github.com/charmbracelet/huh"

func Confirm(title string) bool {
	value := false
	huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(title).
				Affirmative("Да").
				Negative("Нет").
				Value(&value),
		),
	).Run()

	return value
}
