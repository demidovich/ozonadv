package helpers

import (
	"fmt"
	"ozonadv/internal/ozon"

	"github.com/jedib0t/go-pretty/v6/table"
)

func PrintCampaignsTable(campaigns []ozon.Campaign) {
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
	fmt.Println(tw.Render())
	fmt.Println("Всего кампаний:", len(campaigns))
	fmt.Println("")
}
