package ui

import (
	"errors"
	"fmt"
	"ozonadv/internal/cabinets"
	"ozonadv/internal/models"
	"ozonadv/internal/stats"
	"ozonadv/internal/ui/helpers"
	"ozonadv/pkg/console"

	"github.com/charmbracelet/huh"
	"github.com/google/uuid"
)

type cabinetsPage struct {
	cabsService  *cabinets.Service
	statsService *stats.Service
}

func newCabinets(cabsService *cabinets.Service, statsService *stats.Service) cabinetsPage {
	return cabinetsPage{
		cabsService:  cabsService,
		statsService: statsService,
	}
}

func (c cabinetsPage) Run() error {
	options := []helpers.ListOption{}
	for _, cabinet := range c.cabsService.All() {
		options = append(options, helpers.ListOption{
			Key:   cabinet.Name,
			Value: cabinet.UUID,
		})
	}

	options = append(
		options,
		helpers.ListOption{Key: "Новый кабинет", Value: "create_cabinet"},
		helpers.ListOption{Key: "Назад", Value: "back"},
	)

	action, err := helpers.List("---", options...)
	if err != nil {
		return err
	}

	var cabinet models.Cabinet
	var ok bool

	if action == "create_cabinet" {
		cabinet, err = c.createCabinet()
	} else if action == "back" {
		return ErrMainMenu
	} else {
		cabinetUUID := action
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
	helpers.PrintCabinet(cabinet)

	options := []helpers.ListOption{
		{Key: "Кампании", Value: "campaigns_list"},
		{Key: "Отчеты", Value: "stats_list"},
		{Key: "Новый отчет", Value: "create_stat"},
		{Key: "Удалить кабинет", Value: "remove_cabinet"},
		{Key: "Назад", Value: "back"},
	}

	action, err := helpers.List("---", options...)
	if err != nil {
		return err
	}

	switch action {
	case "campaigns_list":
		campaigns, err := c.cabsService.Campaigns(cabinet)
		if err == nil {
			helpers.PrintCampaigns(campaigns)
		}
	case "create_stat":
		_, err = c.createStat(cabinet)
	case "remove_cabinet":
		if console.Confirm("Удалить кабинет \"" + cabinet.Name + "\"?") {
			c.cabsService.Remove(cabinet)
		} else {
			err = c.Run()
		}
	case "back":
		err = c.Run()
	}

	return err
}

func (c cabinetsPage) createCabinet() (models.Cabinet, error) {
	fmt.Println("")
	fmt.Println("Создание нового кабинета")

	cabinet := models.Cabinet{
		UUID: uuid.New().String(),
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Название кабинета").
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

func (c cabinetsPage) createStat(cabinet models.Cabinet) (*models.Stat, error) {
	fmt.Println("Создание нового отчета")

	options := models.StatOptions{}
	if err := helpers.StatOptionsForm(&options); err != nil {
		return nil, err
	}

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
	helpers.PrintStat(*stat)

	return stat, err
}
