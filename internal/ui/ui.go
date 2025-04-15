package ui

import (
	"errors"
	"os"
	"ozonadv/internal/app"
	"ozonadv/internal/ui/helpers"
)

var ErrGoHome = errors.New("go home")
var ErrGoBack = errors.New("go back")

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

	statsPage := newStats(a.StatsService())
	cabinetsPage := newCabinets(a.CabinetsService(), a.StatsService(), statsPage)

	switch action {
	case "cabinets":
		err = cabinetsPage.Home()
	case "stats":
		err = statsPage.Home()
	case "quit":
		os.Exit(0)
	}

	if err != nil && errors.Is(err, ErrGoHome) {
		return Home(a)
	}

	return err
}
