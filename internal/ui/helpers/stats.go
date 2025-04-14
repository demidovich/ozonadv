package helpers

import (
	"fmt"
	"ozonadv/internal/models"

	"github.com/jedib0t/go-pretty/v6/table"
)

func PrintStat(stat models.Stat) {
	tw := table.NewWriter()
	tw.SetStyle(table.StyleRounded)
	tw.AppendRow(table.Row{"Отчет", stat.Options.Name})
	tw.AppendRow(table.Row{"Кабинет", stat.Options.CabinetName})
	tw.AppendRow(table.Row{"Начало интервала, дата", stat.Options.DateFrom})
	tw.AppendRow(table.Row{"Конец интервала, дата", stat.Options.DateTo})
	tw.AppendRow(table.Row{"Группировка", stat.Options.GroupBy})
	tw.AppendRow(table.Row{"Кампаний", len(stat.Items)})

	fmt.Println("")
	fmt.Println(tw.Render())
}
