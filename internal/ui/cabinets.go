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
	"time"

	"github.com/charmbracelet/huh"
	"github.com/google/uuid"
)

type cabinetsPage struct {
	cabsService  *cabinets.Service
	statsService *stats.Service
	statsPage    statsPage
}

func newCabinets(cabsService *cabinets.Service, statsService *stats.Service, statsPage statsPage) cabinetsPage {
	return cabinetsPage{
		cabsService:  cabsService,
		statsService: statsService,
		statsPage:    statsPage,
	}
}

func (c cabinetsPage) Home() error {
	options := []helpers.ListOption{}
	for _, cabinet := range c.cabsService.All() {
		options = append(options, helpers.ListOption{
			Key:   cabinet.Name + " " + helpers.TextGray("(", cabinet.ClientID, ")"),
			Value: cabinet.UUID,
		})
	}

	options = append(
		options,
		helpers.ListOption{Key: "Добавить кабинет", Value: "create_cabinet"},
		helpers.ListOption{Key: "Назад", Value: "back"},
	)

	fmt.Println("")
	action, err := helpers.List("Кабинеты", options...)
	if err != nil {
		return err
	}

	var cabinet *models.Cabinet
	var ok bool

	if action == "create_cabinet" {
		cabinet, err = c.createCabinet()
		if isFormCancel(err) {
			return c.Home()
		}
	} else if action == "back" {
		return ErrGoBack
	} else {
		cabinetUUID := action
		if cabinet, ok = c.cabsService.Find(cabinetUUID); !ok {
			err = errors.New("кабинет не найден")
		}
	}

	if err != nil {
		return err
	}

	return c.cabinet(*cabinet)
}

func (c cabinetsPage) cabinet(cabinet models.Cabinet) error {
	helpers.PrintCabinet(cabinet)

	options := []helpers.ListOption{
		{Key: "Кампании", Value: "campaigns_list"},
		{Key: "Отчеты", Value: "stats_list"},
		{Key: "Редактировать", Value: "update_cabinet"},
		{Key: "Удалить", Value: "remove_cabinet"},
		{Key: "Назад", Value: "back"},
	}

	action, err := helpers.List("Кабинеты > "+cabinet.Name, options...)
	if err != nil {
		return err
	}

	switch action {
	case "campaigns_list":
		campaigns, err := c.cabsService.Campaigns(cabinet)
		if err == nil {
			helpers.PrintCampaigns(campaigns)
			return c.cabinet(cabinet)
		}
	case "stats_list":
		err = c.CabinetStats(cabinet)
	case "update_cabinet":
		err = c.updateCabinet(&cabinet)
		if err == nil || isFormCancel(err) {
			return c.cabinet(cabinet)
		}
	case "remove_cabinet":
		if console.Confirm("Удалить кабинет \"" + cabinet.Name + "\"?") {
			c.cabsService.Remove(cabinet)
			err = ErrGoBack
		} else {
			err = c.cabinet(cabinet)
		}
	case "back":
		err = ErrGoBack
	}

	return err
}

func (c cabinetsPage) createCabinet() (*models.Cabinet, error) {
	cabinet := models.Cabinet{
		UUID:      uuid.New().String(),
		Name:      "Новый кабинет",
		CreatedAt: time.Now().String(),
	}

	err := c.updateCabinet(&cabinet)

	return &cabinet, err
}

func (c cabinetsPage) updateCabinet(cabinet *models.Cabinet) error {
	fmt.Println("Кабинеты > " + cabinet.Name)

	confirm := false
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(forms.RequiredTitle("Название кабинета")).
				CharLimit(100).
				Validate(validators.Required).
				Value(&cabinet.Name),
			huh.NewInput().
				Title(forms.RequiredTitle("Клиент ID")).
				CharLimit(500).
				Validate(validators.Required).
				Value(&cabinet.ClientID),
			huh.NewInput().
				Title(forms.RequiredTitle("Клиент Secret")).
				CharLimit(500).
				Validate(validators.Required).
				Value(&cabinet.ClientSecret),
			huh.NewConfirm().
				Key("done").
				Value(&confirm).
				Inline(true).
				Affirmative("Сохранить").
				Negative("Отмена"),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	if !confirm {
		return ErrFormCancel
	}

	return c.cabsService.Add(*cabinet)
}
