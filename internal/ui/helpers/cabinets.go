package helpers

import (
	"fmt"
	"ozonadv/internal/models"

	"github.com/jedib0t/go-pretty/v6/table"
)

func PrintCabinet(cabinet models.Cabinet) {
	tw := table.NewWriter()
	tw.SetStyle(table.StyleRounded)
	tw.AppendRow(table.Row{"Кабинет", cabinet.Name})
	tw.AppendRow(table.Row{"Клиент ID", cabinet.ClientID})
	tw.AppendRow(table.Row{"Клиент Secret", cabinet.ClientSecretMasked(25)})

	fmt.Println(tw.Render())
	fmt.Println("")
}
