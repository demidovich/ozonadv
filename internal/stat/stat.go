// Ограничения API
// Лимит на количество дней в выгрузке                          62
// Лимит на количество кампаний в отчёте                        10
// Лимит на количество одновременных выгрузок с аккаунта 	    1
// Лимит на количество выгрузок за 24 часа с аккаунта 	        2000
// Лимиты на количество одновременных выгрузок по организации 	5
// Лимит на количество выгрузок за 24 часа в рамках организации 2000

package stat

import (
	"errors"
	"ozonadv/internal/ozon"
	"ozonadv/internal/storage"
	"ozonadv/pkg/validation"
)

type StatOptions struct {
	DateFrom   string `validate:"required,datetime=2006-01-02"`
	DateTo     string `validate:"required,datetime=2006-01-02"`
	ExportFile string `validate:"required,filepath"`
	GroupBy    string `validate:"required,oneof=NO_GROUP_BY,DATE,START_OF_WEEK,START_OF_MONTH"`
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
	s.initReport(options)
	s.processReport(options)

	return nil
}

func (s *statUsecase) HandleContinue() error {
	if s.storage.CampaignRequestsSize() == 0 {
		return errors.New("Кампании для формирования отчета отсутствуют")
	}

	options := StatOptions{
		DateFrom:   s.storage.RequestOptions().DateFrom,
		DateTo:     s.storage.RequestOptions().DateTo,
		ExportFile: s.storage.RequestOptions().ExportFile,
		GroupBy:    s.storage.RequestOptions().GroupBy,
	}

	s.processReport(options)

	return nil
}

func (s *statUsecase) initReport(options StatOptions) {
	storageOptions := storage.RequestOptions{
		DateFrom:   options.DateFrom,
		DateTo:     options.DateTo,
		ExportFile: options.ExportFile,
		GroupBy:    options.GroupBy,
	}

	s.storage.SetRequestOptions(storageOptions)
}

func (s *statUsecase) processReport(options StatOptions) {

}
