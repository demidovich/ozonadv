package stat

import (
	"fmt"
	"ozonadv/internal/ozon"

	"github.com/jedib0t/go-pretty/v6/table"
)

func printCampaignsTable(campaigns []ozon.Campaign) {
	tw := table.NewWriter()
	tw.SetStyle(table.StyleRounded)
	tw.AppendRow(table.Row{"#", "Тип", "Название", "Запуск", "Окончание", "Состояние отчета"})
	tw.AppendRow(table.Row{"", "", "", "", ""})

	for _, c := range campaigns {
		tw.AppendRow(table.Row{
			c.ID,
			c.AdvObjectType,
			// c.StateShort(),
			c.TitleTruncated(45),
			c.FromDate,
			c.ToDate,
			c.StatReportState(),
			// c.Title,
		})
	}

	fmt.Println(tw.Render())
	fmt.Println("Всего:", len(campaigns))
}
