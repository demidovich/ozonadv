package helpers

import (
	"fmt"
	"ozonadv/internal/models"

	"github.com/charmbracelet/huh"
)

func StatOptionsForm(options *models.StatOptions) error {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Название отчета").
				Description("Чтобы было проще понимать что за статистика").
				CharLimit(100).
				Value(&options.Name),
			huh.NewSelect[string]().
				Title("Тип статистики").
				Options(
					huh.NewOption("По рекламным объектам", "OBJECT"),
					huh.NewOption("По рекламным кампаниями", "TOTAL"),
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

	fmt.Println("")
	fmt.Println("Параметры отчета")

	return form.Run()
}
