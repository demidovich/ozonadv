package ui

import (
	"fmt"
	"ozonadv/internal/cabinets"
	"ozonadv/internal/models"

	"github.com/charmbracelet/huh"
	"github.com/jedib0t/go-pretty/v6/table"
)

func printCampaignsTable(campaigns []models.Campaign) {
	tw := table.NewWriter()
	tw.SetStyle(table.StyleRounded)
	tw.AppendRow(table.Row{"#", "Тип", "Кампания", "Запуск", "Окончание", "Статус"})
	tw.AppendRow(table.Row{"", "", "", "", "", ""})

	for i := range campaigns {
		c := &campaigns[i]
		tw.AppendRow(table.Row{
			c.ID,
			c.AdvObjectType,
			c.TitleTruncated(70),
			c.FromDate,
			c.ToDate,
			c.StateShort(),
		})
	}

	fmt.Println("")
	fmt.Println("Рекламные кампании")
	fmt.Println(tw.Render())
	fmt.Println("Всего кампаний:", len(campaigns))
}

func chooseCampaignsForm(cabsService cabinets.Service, cabinet models.Cabinet) ([]models.Campaign, error) {
	fmt.Println("")
	fmt.Println("Выбор рекламных кампаний")

	filters := cabinets.CampaignFilters{}
	confirm := false
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Часть названия кампании или ее ID").
				Description("Можно не заполнять").
				CharLimit(500).
				Value(&filters.Title),
			huh.NewMultiSelect[string]().
				Title("Статусы").
				Description("Один или несколько статусов, можно не заполнять").
				Options(
					huh.NewOption("RUNNING", "CAMPAIGN_STATE_RUNNING"),
					huh.NewOption("PLANNED", "CAMPAIGN_STATE_PLANNED"),
					huh.NewOption("STOPPED", "CAMPAIGN_STATE_STOPPED"),
					huh.NewOption("INACTIVE", "CAMPAIGN_STATE_INACTIVE"),
					huh.NewOption("ARCHIVED", "CAMPAIGN_STATE_ARCHIVED"),
					huh.NewOption("DRAFT", "CAMPAIGN_STATE_MODERATION_DRAFT"),
					huh.NewOption("IN PROGRESS", "CAMPAIGN_STATE_MODERATION_IN_PROGRESS"),
					huh.NewOption("FAILED", "CAMPAIGN_STATE_MODERATION_FAILED"),
					huh.NewOption("FINISHED", "CAMPAIGN_STATE_FINISHED"),
				).
				Value(&filters.States),
			huh.NewConfirm().
				Key("done").
				Value(&confirm).
				Inline(true).
				Affirmative("Ок").
				Negative("Отмена"),
		),
	)

	if err := form.Run(); err != nil {
		return []models.Campaign{}, err
	}

	if !confirm {
		return []models.Campaign{}, ErrFormCancel
	}

	campaigns, err := cabsService.CampaignsFiltered(cabinet, filters)
	if err != nil {
		return campaigns, err
	}

	if len(campaigns) == 0 {
		fmt.Println("")
		printErrorString("Кампании с такими параметрами не найдены")
		return chooseCampaignsForm(cabsService, cabinet)
	}

	return campaigns, nil
}
