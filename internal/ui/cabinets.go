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
	} else if action == "back" {
		return ErrGoHome
	} else {
		cabinetUUID := action
		if cabinet, ok = c.cabsService.Find(cabinetUUID); !ok {
			err = errors.New("кабинет не найден")
		}
	}

	if err != nil {
		return err
	}

	return c.showCabinet(*cabinet)
}

func (c cabinetsPage) showCabinet(cabinet models.Cabinet) error {
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
		}
	case "stats_list":
		err = c.showCabinetStats(cabinet)
	case "update_cabinet":
		err = c.updateCabinet(&cabinet)
	case "remove_cabinet":
		if console.Confirm("Удалить кабинет \"" + cabinet.Name + "\"?") {
			c.cabsService.Remove(cabinet)
		} else {
			err = c.Home()
		}
	case "back":
		err = c.Home()
	}

	return err
}

func (c cabinetsPage) createCabinet() (*models.Cabinet, error) {
	fmt.Println("Кабинеты > Новый кабинет")

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
		return &cabinet, err
	}

	err := c.cabsService.Add(cabinet)

	return &cabinet, err
}

func (c cabinetsPage) updateCabinet(cabinet *models.Cabinet) error {
	fmt.Println("Кабинеты > " + cabinet.Name)

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
		return err
	}

	return c.cabsService.Add(*cabinet)
}

func (c cabinetsPage) showCabinetStats(cabinet models.Cabinet) error {
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
		stat, err = c.createCabinetStat(cabinet)
	} else if action == "back" {
		return c.showCabinet(cabinet)
	} else {
		statUUID := action
		if stat, ok = c.statsService.Find(statUUID); !ok {
			err = errors.New("отчет не найден")
		}
	}

	if err != nil {
		return err
	}

	err = c.statsPage.ShowStat(stat)
	if errors.Is(err, ErrGoBack) {
		return c.showCabinetStats(cabinet)
	}

	return nil
}

func (c cabinetsPage) createCabinetStat(cabinet models.Cabinet) (*models.Stat, error) {
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
