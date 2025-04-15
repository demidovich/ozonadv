package helpers

import (
	"ozonadv/internal/models"

	"github.com/charmbracelet/huh"
)

func StatOptionsForm(options *models.StatOptions) error {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Название отчета").
				CharLimit(100).
				Value(&options.Name),
			huh.NewSelect[string]().
				Title("Тип статистики").
				Options(
					huh.NewOption("Рекламные объекты", "OBJECT"),
					huh.NewOption("Рекламные кампании", "TOTAL"),
				).
				Value(&options.Type),
			huh.NewInput().
				Title("Начало интервала, дата").
				Placeholder("ГГГГ-ДД-ММ").
				CharLimit(10).
				Value(&options.DateFrom),
			huh.NewInput().
				Title("Конец интервала, дата").
				Placeholder("ГГГГ-ДД-ММ").
				CharLimit(10).
				Value(&options.DateTo),
			huh.NewSelect[string]().
				Title("Группировка").
				Options(
					huh.NewOption("Не группировать", "NO_GROUP_BY"),
					huh.NewOption("День", "DATE"),
					huh.NewOption("Неделя", "START_OF_WEEK"),
					huh.NewOption("Месяц", "START_OF_MONTH"),
				).
				Value(&options.GroupBy),
		),
	)

	return form.Run()
}
