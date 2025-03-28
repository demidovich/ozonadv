// Ограничения API
// Лимит на количество дней в выгрузке                          62
// Лимит на количество кампаний в отчёте                        10
// Лимит на количество одновременных выгрузок с аккаунта 	    1
// Лимит на количество выгрузок за 24 часа с аккаунта 	        2000
// Лимиты на количество одновременных выгрузок по организации 	5
// Лимит на количество выгрузок за 24 часа в рамках организации 2000

package object_stat

import (
	"fmt"
	"log"
	"os"
	"ozonadv/internal/ozon"
	"ozonadv/internal/stat/stat_processor"
	"ozonadv/internal/storage"
	"ozonadv/pkg/console"
	"ozonadv/pkg/validation"
	"time"
)

type StatOptions struct {
	DateFrom         string `validate:"required,datetime=2006-01-02"`
	DateTo           string `validate:"required,datetime=2006-01-02"`
	GroupBy          string `validate:"required,oneof=NO_GROUP_BY DATE START_OF_WEEK START_OF_MONTH"`
	CampaignId       string `validate:"omitempty,numeric"`
	CreatedAt        string `validate:"omitempty,datetime=2006-01-02"`
	StartedAt        string `validate:"omitempty,datetime=2006-01-02"`
	ApiRequestsCount int    `validate:"omitempty,numeric"`
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

	options.CreatedAt = time.Now().String()
	options.StartedAt = time.Now().String()
	s.storeOptions(options)

	campaigns := s.selectCampaigns(options)
	for _, c := range campaigns {
		s.storage.ObjectStatCampaigns().Add(c)
	}

	statProcessor := stat_processor.New(s.ozon, s.storage)
	statProcessor.Start(campaigns)

	return nil
}

func (s *statUsecase) HandleContinue() error {
	campaigns := s.storage.ObjectStatCampaigns().All()

	fmt.Println("Есть незавершенное формирование отчета по рекламным объектам")
	//printCampaignsTable(campaigns)
	fmt.Println("")

	if console.Ask("Продолжить?") == false {
		fmt.Println("")
		os.Exit(0)
	}

	fmt.Println("")

	options := s.storage.ObjectStatOptions()
	options.StartedAt = time.Now().String()
	s.storage.SetObjectStatOptions(*options)

	statProcessor := stat_processor.New(s.ozon, s.storage)
	statProcessor.Start(campaigns)

	return nil
}

func (s *statUsecase) storeOptions(options StatOptions) {
	storageOptions := storage.ObjectStatOptions{
		DateFrom:  options.DateFrom,
		DateTo:    options.DateTo,
		GroupBy:   options.GroupBy,
		CreatedAt: options.CreatedAt,
		StartedAt: options.StartedAt,
	}

	s.storage.SetObjectStatOptions(storageOptions)
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

	fmt.Println("")
	//printCampaignsTable(campaigns)
	fmt.Println("")

	if console.Ask("Продолжить?") == false {
		fmt.Println("")
		os.Exit(0)
	}

	fmt.Println("")
	return campaigns
}
