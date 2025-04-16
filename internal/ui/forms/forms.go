package forms

import "ozonadv/internal/ui/helpers"

func RequiredTitle(title string) string {
	return title + helpers.TextGray(", обязательное")
}
