package ui

import (
	"errors"
	"os"

	"github.com/demidovich/ozonadv/internal/app"
	"github.com/demidovich/ozonadv/internal/ui/helpers"
)

type ui struct {
	statsPage    statsPage
	cabinetsPage cabinetsPage
}

func Home(a *app.Application) error {
	options := []helpers.ListOption{
		{Key: "Кабинеты", Value: "cabinets"},
		{Key: "Отчеты", Value: "stats"},
		{Key: "Выход", Value: "quit"},
	}

	action, err := helpers.List("---", options...)
	if err != nil {
		return err
	}

	ui := &ui{}

	ui.statsPage = newStats(a.CabinetsService(), a.StatsService())
	ui.cabinetsPage = newCabinets(a.CabinetsService(), a.StatsService(), ui)

	switch action {
	case "cabinets":
		err = ui.cabinetsPage.Home()
	case "stats":
		err = ui.statsPage.Home()
	case "quit":
		os.Exit(0)
	}

	if err != nil && errors.Is(err, ErrGoBack) {
		return Home(a)
	}

	return err
}
