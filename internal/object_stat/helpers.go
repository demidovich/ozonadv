package object_stat

import (
	"fmt"
	"ozonadv/internal/ozon"
	"ozonadv/internal/storage"

	"github.com/jedib0t/go-pretty/v6/table"
)

func printReportInfo(s *storage.Storage) {
	fmt.Println("Параметры отчета")
	printOptionsTable(*s.ObjectStatOptions())

	fmt.Println("")
	fmt.Println("Кампании отчета")
	printCampaignsTable(s.ObjectStatCampaigns().All())
}

func printOptionsTable(options storage.ObjectStatOptions) {
	createdAt := string([]rune(options.CreatedAt)[:16:16])

	tw := table.NewWriter()
	tw.SetStyle(table.StyleRounded)
	tw.AppendRow(table.Row{"Интервал", options.DateFrom + " - " + options.DateTo})
	tw.AppendRow(table.Row{"Группировка", options.GroupBy})
	tw.AppendRow(table.Row{"Создан", createdAt})

	fmt.Println(tw.Render())
}

func printCampaignsTable(campaigns []ozon.Campaign) {
	tw := table.NewWriter()
	tw.SetStyle(table.StyleRounded)
	tw.AppendRow(table.Row{"#", "Тип", "Кампания", "Запуск", "Окончание", "Состояние отчета"})
	tw.AppendRow(table.Row{"", "", "", "", ""})

	for _, c := range campaigns {
		tw.AppendRow(table.Row{
			c.ID,
			c.AdvObjectType,
			c.TitleTruncated(45),
			c.FromDate,
			c.ToDate,
			// c.StateShort(),
			c.ObjectStatState(),
			// c.Title,
		})
	}

	fmt.Println(tw.Render())
	fmt.Println("Всего кампаний:", len(campaigns))
}
