package ui

import (
	"errors"
	"fmt"
	"ozonadv/internal/cabinets"
	"ozonadv/internal/models"
	"ozonadv/internal/ui/helpers"

	"github.com/charmbracelet/huh"
	"github.com/google/uuid"
)

type cabinetsPage struct {
	cabsService *cabinets.Service
}

func newCabinets(cabsService *cabinets.Service) cabinetsPage {
	return cabinetsPage{cabsService: cabsService}
}

func (c cabinetsPage) Run() error {
	options := []helpers.ListOption{}
	for _, cabinet := range c.cabsService.All() {
		options = append(options, helpers.ListOption{
			Key:   cabinet.Name,
			Value: cabinet.UUID,
		})
	}

	options = append(options, helpers.ListOption{
		Key:   "Новый кабинет",
		Value: "new",
	})

	cabinetUUID, err := helpers.List("---", options...)
	if err != nil {
		return err
	}

	var cabinet models.Cabinet
	var ok bool

	if cabinetUUID == "new" {
		cabinet, err = c.createCabinet()
	} else {
		if cabinet, ok = c.cabsService.Find(cabinetUUID); !ok {
			err = errors.New("кабинет не найден")
		}
	}

	if err != nil {
		return err
	}

	return c.showCabinet(cabinet)
}

func (c cabinetsPage) showCabinet(cabinet models.Cabinet) error {
	helpers.PrintCabinetInfo(cabinet)

	options := []helpers.ListOption{
		{Key: "Рекламные кампании", Value: "campaigns_list"},
		{Key: "Создать отчет", Value: "create_report"},
		{Key: "Удалить кабинет", Value: "remove_cabinet"},
		{Key: "Главное меню", Value: "main_menu"},
	}

	action, err := helpers.List("---", options...)
	if err != nil {
		return err
	}

	switch action {
	case "campaigns_list":
		campaigns, err := c.cabsService.Campaigns(cabinet)
		if err != nil {
			return err
		}
		helpers.PrintCampaignsTable(campaigns)
	case "create_report":
		fmt.Println("create report")
	case "remove_cabinet":
		c.cabsService.Remove(cabinet)
	case "main_menu":
		err = ViewMainMenu
	}

	return err
}

func (c cabinetsPage) createCabinet() (models.Cabinet, error) {
	cabinet := models.Cabinet{
		UUID: uuid.New().String(),
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Название").
				CharLimit(100).
				Value(&cabinet.Name),
			huh.NewInput().
				Title("Клиент ID").
				CharLimit(500).
				Value(&cabinet.ClientID),
			huh.NewInput().
				Title("Клиент Secret").
				CharLimit(500).
				Value(&cabinet.ClientSecret),
		),
	)

	if err := form.Run(); err != nil {
		return cabinet, err
	}

	err := c.cabsService.Add(cabinet)

	return cabinet, err
}
