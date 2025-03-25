// Ограничения API
// Лимит на количество дней в выгрузке                          62
// Лимит на количество кампаний в отчёте                        10
// Лимит на количество одновременных выгрузок с аккаунта 	    1
// Лимит на количество выгрузок за 24 часа с аккаунта 	        2000
// Лимиты на количество одновременных выгрузок по организации 	5
// Лимит на количество выгрузок за 24 часа в рамках организации 2000

package stat

import (
	"fmt"
	"log"
	"os"
	"ozonadv/internal/ozon"
	"ozonadv/internal/stat/stat_processor"
	"ozonadv/internal/storage"
	"ozonadv/pkg/console"
	"ozonadv/pkg/validation"

	"github.com/jedib0t/go-pretty/v6/table"
)

type StatOptions struct {
	DateFrom   string `validate:"required,datetime=2006-01-02"`
	DateTo     string `validate:"required,datetime=2006-01-02"`
	ExportFile string `validate:"required,filepath"`
	GroupBy    string `validate:"required,oneof=NO_GROUP_BY DATE START_OF_WEEK START_OF_MONTH"`
	CampaignId string `validate:"omitempty,numeric"`
}

func (s *StatOptions) Validate() error {
	return validation.ValidateStruct(s)
}

type statUsecase struct {
	storage *storage.Storage
	ozon    *ozon.Ozon
}

func (s *statUsecase) HandleNew(options StatOptions) error {
	if err := options.Validate(); err != nil {
		return err
	}

	if err := s.storage.Reset(); err != nil {
		return err
	}

	s.storeOptions(options)

	campaigns := s.selectCampaigns(options)
	for _, c := range campaigns {
		s.storage.Campaigns.Add(c)
	}

	s.startPocessing(campaigns)

	return nil
}

func (s *statUsecase) HandleContinue() error {
	campaigns := s.storage.Campaigns.All()

	fmt.Printf("")
	s.printCampaigns(campaigns)
	fmt.Println("")

	if console.Ask("Продолжить?") == false {
		fmt.Println("")
		os.Exit(0)
	}

	s.startPocessing(campaigns)

	return nil
}

func (s *statUsecase) storeOptions(options StatOptions) {
	storageOptions := storage.StatOptions{
		DateFrom:   options.DateFrom,
		DateTo:     options.DateTo,
		ExportFile: options.ExportFile,
		GroupBy:    options.GroupBy,
	}

	s.storage.SetStatOptions(storageOptions)
}

func (s *statUsecase) selectCampaigns(options StatOptions) []ozon.Campaign {
	filters := ozon.FindCampaignsFilters{}
	if options.CampaignId != "" {
		filters.Ids = append(filters.Ids, options.CampaignId)
	}

	campaigns, err := s.ozon.Campaigns().Find(filters)
	if err != nil {
		log.Fatal(err)
	}

	if len(campaigns) == 0 {
		log.Fatal("Кампании не найдены")
	}

	fmt.Printf("")
	s.printCampaigns(campaigns)
	fmt.Println("")

	if console.Ask("Продолжить?") == false {
		fmt.Println("")
		os.Exit(0)
	}

	return campaigns
}

func (s *statUsecase) startPocessing(campaigns []ozon.Campaign) {
	campaignsCh := make(chan ozon.Campaign)

	go func() {
		defer close(campaignsCh)
		for _, campaign := range campaigns {
			campaignsCh <- campaign
		}
	}()

	statFiles := stat_processor.Start(s.ozon, s.storage, campaignsCh)
	for file := range statFiles {
		fmt.Println(file)
	}
}

func (s *statUsecase) printCampaigns(campaigns []ozon.Campaign) {
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
