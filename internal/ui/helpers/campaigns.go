package helpers

import (
	"fmt"
	"ozonadv/internal/cabinets"
	"ozonadv/internal/models"
	"ozonadv/internal/ozon"

	"github.com/charmbracelet/huh"
	"github.com/jedib0t/go-pretty/v6/table"
)

func PrintCampaigns(campaigns []ozon.Campaign) {
	tw := table.NewWriter()
	tw.SetStyle(table.StyleRounded)
	tw.AppendRow(table.Row{"#", "Тип", "Кампания", "Запуск", "Окончание", "Статус"})
	tw.AppendRow(table.Row{"", "", "", "", "", ""})

	for _, c := range campaigns {
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

func ChooseCampaigns(cabsService cabinets.Service, cabinet models.Cabinet) ([]ozon.Campaign, error) {
	fmt.Println("Выбор рекламных кампаний")

	filters := cabinets.CampaignFilters{}
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Часть названия кампании").
				Description("Можно не заполнять").
				CharLimit(100).
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
		),
	)

	if err := form.Run(); err != nil {
		return []ozon.Campaign{}, err
	}

	campaigns, err := cabsService.CampaignsFiltered(cabinet, filters)
	if err != nil {
		return campaigns, err
	}

	if len(campaigns) == 0 {
		fmt.Println("")
		fmt.Println("Кампании с такими параметрами не найдены")
		fmt.Println("")
		return ChooseCampaigns(cabsService, cabinet)
	}

	return campaigns, nil
}
