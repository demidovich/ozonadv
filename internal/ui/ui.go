package ui

import (
	"os"
	"ozonadv/internal/app"
	"ozonadv/pkg/prompts"
)

func Run(a *app.Application) error {

	options := []prompts.ListOption{
		{Key: "Кабинеты", Value: "cabinets"},
		{Key: "Статистика", Value: "stats"},
		{Key: "Выход", Value: "quit"},
	}

	action, err := prompts.List("---", options...)
	if err != nil {
		return err
	}

	switch action {
	case "cabinets":
		err = cabnets(a)
	case "stats":
		err = stats(a)
	case "quit":
		os.Exit(0)
	}

	return err
}
