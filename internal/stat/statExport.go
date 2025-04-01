package stat

import (
	"fmt"
	"ozonadv/internal/ozon"
	"ozonadv/internal/stat/stat_export"
	"ozonadv/internal/storage"
	"ozonadv/pkg/console"
	"ozonadv/pkg/validation"
)

type statExportUsecase struct {
	storage *storage.Storage
}

type StatExportOptions struct {
	// File string `validate:"required,filepath"`
	File string `validate:"required"`
}

func (s *StatExportOptions) Validate() error {
	return validation.ValidateStruct(s)
}

func (s *statExportUsecase) Handle(options StatExportOptions) error {
	if err := options.Validate(); err != nil {
		return err
	}

	if s.storage.StatCampaigns().Size() == 0 {
		fmt.Println("Обработанных кампаний нет")
		return nil
	}

	campaigns := s.storage.StatCampaigns().All()
	printCampaignsTable(campaigns)

	incompleted := s.incompletedCampaigns(campaigns)
	if len(incompleted) > 0 {
		fmt.Println("")
		fmt.Println("Найдены кампании с несформированой статистикой")
		fmt.Println("Экспорт по ним не будет выполнен")

		fmt.Println("")
		printCampaignsTable(incompleted)

		fmt.Println("")
		if console.Ask("Подолжить?") == false {
			return nil
		}
	}

	statExport := stat_export.New(s.storage)
	statExport.ToFile(options.File)

	return nil
}

func (s *statExportUsecase) incompletedCampaigns(campaigns []ozon.Campaign) []ozon.Campaign {
	incompleted := []ozon.Campaign{}
	for _, c := range campaigns {
		if c.Stat.File == "" {
			incompleted = append(incompleted, c)
		}
	}

	return incompleted
}
