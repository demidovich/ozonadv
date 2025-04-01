package find

import (
	"fmt"
	"ozonadv/internal/ozon"
	"ozonadv/internal/storage"
	"ozonadv/pkg/console"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

type findCampaignsUsecase struct {
	storage *storage.Storage
	ozon    *ozon.Ozon
}

func (f *findCampaignsUsecase) Handle() error {
	all, err := f.ozon.Campaigns().All()
	if err != nil {
		return err
	}

	if len(all) == 0 {
		fmt.Println("Кампании не найдены.")
		return nil
	}

	var campaigns []ozon.Campaign

	fmt.Println("")
	if console.Ask("Отфильтровать?") == true {
		campaigns = f.filteredByTitle(all)
		campaigns = f.filteredByState(campaigns)
	} else {
		campaigns = all
	}

	f.printCampaignsTable(campaigns)

	if console.Ask("Использовать кампании для отчета по рекламным объектам?") == false {
		return nil
	}

	f.storage.ObjectStatResetAll()
	for _, c := range campaigns {
		f.storage.ObjectStatCampaigns().Add(c)
	}

	return nil
}

func (f *findCampaignsUsecase) filteredByTitle(all []ozon.Campaign) []ozon.Campaign {
	value := console.InputString("Часть имени")

	var filtered []ozon.Campaign
	for _, c := range all {
		cTitle := strings.ToLower(c.Title)
		iTitle := strings.ToLower(value)
		if strings.Contains(cTitle, iTitle) {
			filtered = append(filtered, c)
		}
	}

	return filtered
}

func (f *findCampaignsUsecase) filteredByState(all []ozon.Campaign) []ozon.Campaign {
	value := console.InputString("Часть статуса")

	var filtered []ozon.Campaign
	for _, c := range all {
		cState := strings.ToLower(c.State)
		iState := strings.ToLower(value)
		if strings.Contains(cState, iState) {
			filtered = append(filtered, c)
		}
	}

	return filtered
}

func (s *findCampaignsUsecase) printCampaignsTable(campaigns []ozon.Campaign) {
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
