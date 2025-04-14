package ui

import (
	"errors"
	"os"
	"ozonadv/internal/app"
	"ozonadv/internal/ui/helpers"
)

var ViewMainMenu = errors.New("main menu")

func Run(a *app.Application) error {
	options := []helpers.ListOption{
		{Key: "Кабинеты", Value: "cabinets"},
		{Key: "Отчеты", Value: "stats"},
		{Key: "Выход", Value: "quit"},
	}

	action, err := helpers.List("---", options...)
	if err != nil {
		return err
	}

	switch action {
	case "cabinets":
		err = newCabinets(a.CabinetsService()).Run()
	case "stats":
		err = newStats(a.StatsService()).Run()
	case "quit":
		os.Exit(0)
	}

	if err != nil && errors.Is(err, ViewMainMenu) {
		return Run(a)
	}

	return err
}
