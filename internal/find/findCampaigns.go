package find

import (
	"fmt"
	"ozonadv/internal/ozon"

	"github.com/jedib0t/go-pretty/v6/table"
)

type findCampaignsUsecase struct {
	ozon *ozon.Ozon
}

func (f *findCampaignsUsecase) Handle() error {
	campaigns, err := f.ozon.Campaigns().All()
	if err != nil {
		return err
	}

	if len(campaigns) == 0 {
		fmt.Println("Кампании не найдены.")
		return nil
	}

	f.printCampaignsTable(campaigns)

	return nil
}

func (s *findCampaignsUsecase) printCampaignsTable(campaigns []ozon.Campaign) {
	tw := table.NewWriter()
	tw.SetStyle(table.StyleRounded)
	tw.AppendRow(table.Row{"#", "State", "Type", "From", "To", "Title"})
	tw.AppendRow(table.Row{"", "", "", "", "", ""})

	for _, c := range campaigns {
		tw.AppendRow(table.Row{
			c.ID,
			c.ShortState(),
			c.AdvObjectType,
			c.FromDate,
			c.ToDate,
			c.Title,
		})
	}

	fmt.Println(tw.Render())
	fmt.Println("Всего:", len(campaigns))
}
