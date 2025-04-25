package ui

import (
	"fmt"

	"github.com/demidovich/ozonadv/internal/cabinets"
	"github.com/demidovich/ozonadv/internal/models"

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
				Description("Предзаполнены статусы кампаний, которые могли работать").
				Options(
					huh.NewOption("RUNNING", "CAMPAIGN_STATE_RUNNING").Selected(true),
					huh.NewOption("INACTIVE", "CAMPAIGN_STATE_INACTIVE").Selected(true),
					huh.NewOption("STOPPED", "CAMPAIGN_STATE_STOPPED").Selected(true),
					huh.NewOption("FINISHED", "CAMPAIGN_STATE_FINISHED").Selected(true),
					huh.NewOption("ARCHIVED", "CAMPAIGN_STATE_ARCHIVED").Selected(true),
					huh.NewOption("MODERATION DRAFT", "CAMPAIGN_STATE_MODERATION_DRAFT").Selected(true),
					huh.NewOption("MODERATION IN PROGRESS", "CAMPAIGN_STATE_MODERATION_IN_PROGRESS").Selected(true),
					huh.NewOption("MODERATION FAILED", "CAMPAIGN_STATE_MODERATION_FAILED").Selected(true),
					huh.NewOption("PLANNED", "CAMPAIGN_STATE_PLANNED"),
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
