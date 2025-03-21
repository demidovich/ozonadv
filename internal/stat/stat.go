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
	"ozonadv/internal/storage"
	"ozonadv/pkg/console"
	"ozonadv/pkg/validation"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/schollz/progressbar/v3"
)

type StatOptions struct {
	DateFrom   string `validate:"required,datetime=2006-01-02"`
	DateTo     string `validate:"required,datetime=2006-01-02"`
	CampaignId string `validate:"omitempty,numeric"`
	ExportFile string `validate:"required,filepath"`
	GroupBy    string `validate:"required,oneof=NO_GROUP_BY DATE START_OF_WEEK START_OF_MONTH"`
}

func (c *StatOptions) Validate() error {
	return validation.ValidateStruct(c)
}

type statUsecase struct {
	storage *storage.Storage
	ozonApi *ozon.Api
}

func (s *statUsecase) HandleNew(options StatOptions) error {
	if err := options.Validate(); err != nil {
		return err
	}

	if err := s.storage.Reset(); err != nil {
		return err
	}

	s.initProcessingOptions(options)
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

func (s *statUsecase) initCampaigns(options StatOptions) error {
	filters := ozon.FindCampaignsFilters{}
	if options.CampaignId != "" {
		filters.Ids = append(filters.Ids, options.CampaignId)
	}

	campaigns, err := s.ozonApi.FindCampaigns(filters)
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

func (s *statUsecase) initProcessingOptions(options StatOptions) {
	storageOptions := storage.RequestOptions{
		DateFrom:   options.DateFrom,
		DateTo:     options.DateTo,
		ExportFile: options.ExportFile,
		GroupBy:    options.GroupBy,
	}

	s.storage.SetRequestOptions(storageOptions)
}

func (s *statUsecase) startPocessing() {
	max := s.storage.Campaigns.Size()
	bar := progressbar.Default(int64(max))

	for {
		campaign, ok := s.storage.Campaigns.Next()
		if !ok {
			break
		}

		s.processCampaign(campaign)
		bar.Add(1)
	}
}

func (s *statUsecase) processCampaign(campaign ozon.Campaign) error {
	// fmt.Printf("#%s, %s", campaign.ID, campaign.Title)

	time.Sleep(5 * time.Second)

	s.storage.Campaigns.Remove(campaign.ID)
	return nil
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

// func (s *statUsecase) printCampaigns(campaigns map[string]ozon.Campaign) {
// 	statMaxLen := 0
// 	typeMaxLen := 0

// 	for _, c := range campaigns {
// 		statLen := utf8.RuneCountInString(c.ShortState())
// 		if statLen > statMaxLen {
// 			statMaxLen = statLen
// 		}
// 		typeLen := utf8.RuneCountInString(c.AdvObjectType)
// 		if typeLen > typeMaxLen {
// 			typeMaxLen = typeLen
// 		}
// 	}

// 	format := "#%-9s %-" + strconv.Itoa(statMaxLen) + "s  %-" + strconv.Itoa(typeMaxLen) + "s  %s\n"
// 	for _, c := range campaigns {
// 		fmt.Printf(format, c.ID, c.ShortState(), c.AdvObjectType, c.Title)
// 	}

// 	fmt.Println("Всего:", len(campaigns))
// }
