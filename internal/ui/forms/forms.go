package forms

import "ozonadv/internal/ui/colors"

func RequiredTitle(title string) string {
	return title + colors.Gray().Sprint(", обязательное")
}
