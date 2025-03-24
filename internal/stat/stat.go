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

	s.initOptions(options)

	if err := s.initCampaigns(options); err != nil {
		return err
	}

	s.startPocessing()

	return nil
}

func (s *statUsecase) HandleContinue() error {
	fmt.Printf("")
	s.printCampaigns(s.storage.Campaigns.All())
	fmt.Println("")

	s.startPocessing()

	return nil
}

func (s *statUsecase) initOptions(options StatOptions) {
	storageOptions := storage.StatOptions{
		DateFrom:   options.DateFrom,
		DateTo:     options.DateTo,
		ExportFile: options.ExportFile,
		GroupBy:    options.GroupBy,
	}

	s.storage.SetStatOptions(storageOptions)
}

func (s *statUsecase) initCampaigns(options StatOptions) error {
	filters := ozon.FindCampaignsFilters{}
	if options.CampaignId != "" {
		filters.Ids = append(filters.Ids, options.CampaignId)
	}

	campaigns, err := s.ozon.Campaigns().Find(filters)
	if err != nil {
		return err
	}

	for _, c := range campaigns {
		if c.NeverRun() {
			continue
		}
		s.storage.Campaigns.Add(c)
	}

	if s.storage.Campaigns.Size() == 0 {
		fmt.Println("Кампании, которые могли работать, не найдены")
		fmt.Println("")
		os.Exit(0)
	}

	fmt.Printf("")
	s.printCampaigns(s.storage.Campaigns.All())
	fmt.Println("")

	if console.Ask("Продолжить?") == false {
		fmt.Println("")
		os.Exit(0)
	}
	fmt.Println("")

	return nil
}

func (s *statUsecase) startPocessing() {
	campaigns := make(chan ozon.Campaign)
	defer close(campaigns)
	go func() {
		for _, campaign := range s.storage.Campaigns.All() {
			campaigns <- campaign
		}
	}()

	statFiles := stat_processor.Start(s.ozon, s.storage, campaigns)
	for file := range statFiles {
		fmt.Println(file)
		fmt.Println("222222222222222")
	}

	fmt.Println("33333333333")
}

func (s *statUsecase) printCampaigns(campaigns map[string]ozon.Campaign) {
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
