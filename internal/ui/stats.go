package ui

import (
	"fmt"
	"ozonadv/internal/models"
	"ozonadv/internal/stats"
	"ozonadv/internal/ui/helpers"
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

func (c statsPage) ShowStat(stat *models.Stat) error {
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
