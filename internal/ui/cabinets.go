package ui

import (
	"errors"
	"fmt"
	"ozonadv/internal/cabinets"
	"ozonadv/internal/models"
	"ozonadv/internal/stats"
	"ozonadv/internal/ui/colors"
	"ozonadv/internal/ui/forms"
	"ozonadv/internal/ui/forms/validators"
	"ozonadv/internal/ui/helpers"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/google/uuid"
	"github.com/jedib0t/go-pretty/v6/table"
)

type cabinetsPage struct {
	cabsService  *cabinets.Service
	statsService *stats.Service
	ui           *ui
}

func newCabinets(cabsService *cabinets.Service, statsService *stats.Service, ui *ui) cabinetsPage {
	return cabinetsPage{
		cabsService:  cabsService,
		statsService: statsService,
		ui:           ui,
	}
}

func (c cabinetsPage) Home() error {
	options := []helpers.ListOption{}
	for _, cabinet := range c.cabsService.All() {
		options = append(options, helpers.ListOption{
			Key:   cabinet.Name + " " + colors.Gray("(%s)", cabinet.ClientID),
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
		if err == nil || isFormCanceled(err) {
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
	c.printCabinetTable(cabinet)

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
			printCampaignsTable(campaigns)
			fmt.Println("")
			return c.cabinet(cabinet)
		}
	case "stats_list":
		err = c.ui.statsPage.CabinetStats(cabinet)
	case "update_cabinet":
		err = c.updateCabinet(&cabinet)
		if err == nil || isFormCanceled(err) {
			return c.cabinet(cabinet)
		}
	case "remove_cabinet":
		if helpers.Confirm("Удалить кабинет \"" + cabinet.Name + "\"?") {
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
	fmt.Println("Кабинеты > Новый кабинет")

	cabinet := &models.Cabinet{
		UUID:      uuid.New().String(),
		CreatedAt: time.Now().String(),
	}

	if err := c.editCabinetForm(cabinet); err != nil {
		return nil, err
	}

	if err := c.cabsService.Add(*cabinet); err != nil {
		return nil, err
	}

	return cabinet, nil
}

func (c cabinetsPage) updateCabinet(cabinet *models.Cabinet) error {
	fmt.Println("Кабинеты > " + cabinet.Name)

	if err := c.editCabinetForm(cabinet); err != nil {
		return err
	}

	if err := c.cabsService.Add(*cabinet); err != nil {
		return err
	}

	return nil
}

func (c cabinetsPage) editCabinetForm(cabinet *models.Cabinet) error {
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

	return nil
}

func (c cabinetsPage) printCabinetTable(cabinet models.Cabinet) {
	tw := table.NewWriter()
	tw.SetStyle(table.StyleRounded)
	tw.AppendRow(table.Row{"Кабинет", cabinet.Name})
	tw.AppendRow(table.Row{"Клиент ID", cabinet.ClientID})
	tw.AppendRow(table.Row{"Клиент Secret", cabinet.ClientSecretMasked(25)})

	fmt.Println(tw.Render())
	fmt.Println("")
}
