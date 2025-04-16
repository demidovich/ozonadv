package ui

import (
	"errors"
	"fmt"
	"ozonadv/internal/cabinets"
	"ozonadv/internal/models"
	"ozonadv/internal/stats"
	"ozonadv/internal/ui/forms"
	"ozonadv/internal/ui/forms/validators"
	"ozonadv/internal/ui/helpers"
	"ozonadv/pkg/console"

	"github.com/charmbracelet/huh"
	"github.com/jedib0t/go-pretty/v6/table"
)

type statsPage struct {
	cabsService  *cabinets.Service
	statsService *stats.Service
}

func newStats(cabsService *cabinets.Service, statsService *stats.Service) statsPage {
	return statsPage{
		cabsService:  cabsService,
		statsService: statsService,
	}
}

func (c statsPage) Home() error {
	return nil
}

func (c statsPage) Stat(stat *models.Stat) error {
	c.printStatTable(*stat)
	fmt.Println("")

	options := []helpers.ListOption{
		{Key: "Загрузка", Value: "download"},
		{Key: "Экспорт", Value: "export"},
		{Key: "Назад", Value: "back"},
	}

	action, err := helpers.List("Кабинеты > "+stat.Options.CabinetName+" > Отчеты > "+stat.Options.Name, options...)
	if err != nil {
		return err
	}

	if action == "download" {
		c.statsService.Download(stat)
	} else if action == "back" {
		return ErrGoBack
	}

	return nil
}

func (c statsPage) CabinetStats(cabinet models.Cabinet) error {
	options := []helpers.ListOption{}
	for _, stat := range c.statsService.CabinetAll(cabinet) {
		options = append(options, helpers.ListOption{
			Key:   stat.Options.Name,
			Value: stat.UUID,
		})
	}

	options = append(
		options,
		helpers.ListOption{Key: "Добавить отчет", Value: "create_stat"},
		helpers.ListOption{Key: "Назад", Value: "back"},
	)

	action, err := helpers.List("Кабинеты > "+cabinet.Name+" > Отчеты", options...)
	if err != nil {
		return err
	}

	var stat *models.Stat
	var ok bool

	if action == "create_stat" {
		stat, err = c.CreateStat(cabinet)
		if isFormCanceled(err) {
			fmt.Println("")
			return c.CabinetStats(cabinet)
		}
	} else if action == "back" {
		return ErrGoBack
	} else {
		statUUID := action
		if stat, ok = c.statsService.Find(statUUID); !ok {
			err = errors.New("отчет не найден")
		}
	}

	if err != nil {
		return err
	}

	err = c.Stat(stat)
	if errors.Is(err, ErrGoBack) {
		return c.CabinetStats(cabinet)
	}

	return nil
}

func (c statsPage) CreateStat(cabinet models.Cabinet) (*models.Stat, error) {
	fmt.Println(cabinet.Name + " > Новый отчет")

	options := models.StatOptions{}
	if err := c.statOptionsForm(&options); err != nil {
		return nil, err
	}

	options.CabinetUUID = cabinet.UUID
	options.CabinetName = cabinet.Name
	options.CabinetClientId = cabinet.ClientID
	options.CabinetClientSecret = cabinet.ClientSecret

	campaigns, err := chooseCampaignsForm(*c.cabsService, cabinet)

	if isFormCanceled(err) {
		return nil, ErrFormCancel
	}

	if err != nil {
		return nil, err
	}

	printCampaignsTable(campaigns)

	fmt.Println("")
	if !console.Confirm("Создать отчет?") {
		return nil, ErrFormCancel
	}

	stat, err := c.statsService.Create(options, campaigns)
	if err != nil {
		return stat, err
	}

	fmt.Println("")
	fmt.Println("Отчет создан")
	fmt.Println("")

	return stat, err
}

func (c statsPage) statOptionsForm(options *models.StatOptions) error {
	confirm := false
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(forms.RequiredTitle("Название отчета")).
				CharLimit(100).
				Validate(validators.Required).
				Value(&options.Name),
			huh.NewSelect[string]().
				Title(forms.RequiredTitle("Тип статистики")).
				Options(
					huh.NewOption("Рекламные объекты", "OBJECT"),
					huh.NewOption("Рекламные кампании", "TOTAL"),
				).
				Validate(validators.Required).
				Value(&options.Type),
			huh.NewInput().
				Title(forms.RequiredTitle("Начало интервала, дата")).
				Placeholder("ГГГГ-ДД-ММ").
				CharLimit(10).
				Validate(validators.DateRequiured).
				Value(&options.DateFrom),
			huh.NewInput().
				Title(forms.RequiredTitle("Конец интервала, дата")).
				Placeholder("ГГГГ-ДД-ММ").
				CharLimit(10).
				Validate(validators.DateRequiured).
				Value(&options.DateTo),
			huh.NewSelect[string]().
				Title(forms.RequiredTitle("Группировка")).
				Options(
					huh.NewOption("Не группировать", "NO_GROUP_BY"),
					huh.NewOption("День", "DATE"),
					huh.NewOption("Неделя", "START_OF_WEEK"),
					huh.NewOption("Месяц", "START_OF_MONTH"),
				).
				Validate(validators.Required).
				Value(&options.GroupBy),
			huh.NewConfirm().
				Key("done").
				Value(&confirm).
				Inline(true).
				Negative("Отмена").
				Affirmative("Далее"),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	if !confirm {
		return ErrFormCancel
	}

	return nil
}

func (c statsPage) printStatTable(stat models.Stat) {
	tw := table.NewWriter()
	tw.SetStyle(table.StyleRounded)
	tw.AppendRow(table.Row{"Отчет", stat.Options.Name})
	tw.AppendRow(table.Row{"Кабинет", stat.Options.CabinetName})
	tw.AppendRow(table.Row{"Начало интервала, дата", stat.Options.DateFrom})
	tw.AppendRow(table.Row{"Конец интервала, дата", stat.Options.DateTo})
	tw.AppendRow(table.Row{"Группировка", stat.Options.GroupBy})
	tw.AppendRow(table.Row{"Кампаний", len(stat.Items)})
	tw.AppendRow(table.Row{"Состояние", stat.StateHuman()})

	fmt.Println(tw.Render())
}
