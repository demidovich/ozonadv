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
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/schollz/progressbar/v3"
)

type StatOptions struct {
	DateFrom   string `validate:"required,datetime=2006-01-02"`
	DateTo     string `validate:"required,datetime=2006-01-02"`
	ExportFile string `validate:"required,filepath"`
	GroupBy    string `validate:"required,oneof=NO_GROUP_BY DATE START_OF_WEEK START_OF_MONTH"`
}

func (c *StatOptions) Validate() error {
	return validation.ValidateStruct(c)
}

type statUsecase struct {
	storage    *storage.Storage
	ozonClient *ozon.Client
}

func (s *statUsecase) HandleNew(options StatOptions) error {
	if err := options.Validate(); err != nil {
		return err
	}

	s.storage.Reset()
	s.initProcessingOptions(options)
	if err := s.initCampaigns(); err != nil {
		return err
	}

	s.startPocessing()

	return nil
}

func (s *statUsecase) HandleContinue() error {
	fmt.Printf("")
	s.printCampaignRequests(s.storage.CampaignRequests())
	fmt.Println("")

	s.startPocessing()

	return nil
}

func (s *statUsecase) initCampaigns() error {
	campaigns, err := s.ozonClient.AllCampaigns()
	if err != nil {
		return err
	}

	for _, c := range campaigns {
		if c.NeverRun() {
			continue
		}
		s.storage.AddCampaignRequest(c)
	}

	if s.storage.CampaignRequestsSize() == 0 {
		fmt.Println("Кампании, которые могли работать, не найдены")
		fmt.Println("")
		os.Exit(0)
	}

	fmt.Printf("")
	s.printCampaignRequests(s.storage.CampaignRequests())
	fmt.Println("")

	if console.Ask("Продолжить?") == false {
		s.storage.Reset()
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
	max := s.storage.CampaignRequestsSize()
	bar := progressbar.Default(int64(max))

	for {
		campaign, ok := s.storage.NextCampaignRequest()
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

	s.storage.RemoveCampaignRequest(campaign.ID)
	return nil
}

func (s *statUsecase) printCampaignRequests(campaigns []ozon.Campaign) {
	statMaxLength := 0
	typeMaxLength := 0

	for _, c := range campaigns {
		statLength := utf8.RuneCountInString(c.ShortState())
		if statLength > statMaxLength {
			statMaxLength = statLength
		}
		typeLength := utf8.RuneCountInString(c.AdvObjectType)
		if typeLength > typeMaxLength {
			typeMaxLength = typeLength
		}
	}

	format := "#%-9s %-" + strconv.Itoa(statMaxLength) + "s  %-" + strconv.Itoa(typeMaxLength) + "s  %s\n"
	for _, c := range campaigns {
		fmt.Printf(format, c.ID, c.ShortState(), c.AdvObjectType, c.Title)
	}

	fmt.Println("Всего:", len(campaigns))
}
