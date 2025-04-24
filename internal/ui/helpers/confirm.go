package helpers

import (
	"log"

	"github.com/charmbracelet/huh"
)

func Confirm(title string) bool {
	value := false

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(title).
				Affirmative("Да").
				Negative("Нет").
				Value(&value),
		),
	).Run()

	if err != nil {
		log.Fatal(err)
	}

	return value
}

func WaitButton(text string) {
	_ = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Negative("").
				Affirmative(text),
		),
	).Run()
}
