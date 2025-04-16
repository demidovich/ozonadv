package ui

import (
	"errors"
	"fmt"
	"ozonadv/internal/models"
	"ozonadv/internal/stats"
	"ozonadv/internal/ui/helpers"
	"ozonadv/pkg/console"
)

type statsPage struct {
	statsService *stats.Service
}

func newStats(statsService *stats.Service) statsPage {
	return statsPage{statsService: statsService}
}

func (c statsPage) Home() error {
	return nil
}

func (c statsPage) Stat(stat *models.Stat) error {
	helpers.PrintStat(*stat)
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

func (c cabinetsPage) CabinetStats(cabinet models.Cabinet) error {
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

	err = c.statsPage.Stat(stat)
	if errors.Is(err, ErrGoBack) {
		return c.CabinetStats(cabinet)
	}

	return nil
}

func (c cabinetsPage) CreateStat(cabinet models.Cabinet) (*models.Stat, error) {
	fmt.Println(cabinet.Name + " > Новый отчет")

	options := models.StatOptions{}
	if err := helpers.StatOptionsForm(&options); err != nil {
		return nil, err
	}

	options.CabinetUUID = cabinet.UUID
	options.CabinetName = cabinet.Name
	options.CabinetClientId = cabinet.ClientID
	options.CabinetClientSecret = cabinet.ClientSecret

	campaigns, err := helpers.ChooseCampaigns(*c.cabsService, cabinet)
	if err != nil {
		return nil, err
	}

	helpers.PrintCampaigns(campaigns)

	fmt.Println("")
	if !console.Confirm("Создать отчет?") {
		return nil, errors.New("cancel")
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
