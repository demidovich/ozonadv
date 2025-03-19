package stat

import (
	"fmt"
	"os"
	"ozonadv/internal/ozon"
	"ozonadv/internal/storage"

	"github.com/go-playground/validator/v10"
	"github.com/go-resty/resty/v2"
)

type PullOptions struct {
	ExportFile string `validate:"required,filepath"`
}

func (p *PullOptions) validate() error {
	validate := validator.New()

	err := validate.Struct(p)
	if err == nil {
		return nil
	}

	errs := err.(validator.ValidationErrors)
	if errs != nil {
		return fmt.Errorf("%s", errs)
	}

	return nil
}

type pullUsecase struct {
	storage    *storage.Storage
	ozonClient *ozon.Client
	resty      *resty.Client
}

func (p *pullUsecase) Handle(options PullOptions) error {
	if err := options.validate(); err != nil {
		return err
	}

	if p.storage.StatisticsSize() == 0 {
		fmt.Println("Ожидающие обработки запросы отсутствуют")
		return nil
	}

	ozonStats, err := p.ozonClient.Statistics()
	if err != nil {
		return err
	}

	if len(ozonStats) == 0 {
		fmt.Println("В Озон нет отчетов для обработки")
		return nil
	}

	export, err := os.Open(options.ExportFile)
	if err != nil {
		return err
	}
	defer export.Close()

	for _, s := range ozonStats {
		if !s.IsReadyToDownload() {
			continue
		}

		if !p.storage.HasStatistic(s.UUID) {
			continue
		}

		data, err := p.ozonClient.DownloadStatistic(s.Link)
		if err != nil {
			fmt.Printf("[%s] Download error: %v\n", s.UUID, err)
			continue
		}

		_, err = export.Write(data)
		if err != nil {
			fmt.Printf("[%s] Write error: %v\n", s.UUID, err)
		}
	}

	return nil
}
