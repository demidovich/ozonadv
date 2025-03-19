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
	"ozonadv/internal/ozon"
	"ozonadv/internal/storage"

	"github.com/go-playground/validator/v10"
)

type CreateOptions struct {
	FromDate string `validate:"required,datetime=2006-01-02"`
	ToDate   string `validator:"required,datetime=2006-01-02"`
}

func (c *CreateOptions) validate() error {
	validate := validator.New()

	err := validate.Struct(c)
	if err == nil {
		return nil
	}

	errs := err.(validator.ValidationErrors)
	return fmt.Errorf("%s", errs)
}

type createUsecase struct {
	storage    *storage.Storage
	ozonClient *ozon.Client
}

func (c *createUsecase) Handle(options CreateOptions) error {
	if err := options.validate(); err != nil {
		return err
	}

	campaigns, err := c.ozonClient.Campaigns()
	if err != nil {
		return err
	}

	fmt.Println("")
	fmt.Printf("Найдено кампаний: %d\n", len(campaigns))
	fmt.Println("")

	for _, campaign := range campaigns {
		fmt.Printf("#%s, %s, %s\n", campaign.ID, campaign.State, campaign.Title)

		if campaign.NotRunning() {
			fmt.Println("Пропуск. Кампания не была запущена.")
			fmt.Println("")
			continue
		}

		err := c.createStatistic(campaign, options)
		if err != nil {
			fmt.Println("Ошибка запроса формирования отчета:", err)
			continue
		}

		fmt.Println("Отправлен запрос формирования отчета.")
		fmt.Println("")
	}

	return nil
}

func (c *createUsecase) createStatistic(campaign ozon.Campaign, options CreateOptions) error {
	var resource string
	var payload map[string]any

	if campaign.AdvObjectType == "VIDEO_BANNER" {
		resource = "/api/adv-api/external/api/statistics/video"
		payload = map[string]any{
			"campaigns": []string{campaign.ID},
			"dateFrom":  options.FromDate,
			"dateTo":    options.ToDate,
			"groupBy":   "DAY",
		}
	} else {
		resource = "/api/adv-api/external/api/statistics"
		payload = map[string]any{
			"campaignId": campaign.ID,
			"dateFrom":   options.FromDate,
			"dateTo":     options.ToDate,
			"groupBy":    "DAY",
		}
	}

	result := struct {
		UUID   string `json:"UUID"`
		Vendor bool   `json:"vendor"`
	}{}

	err := c.ozonClient.Post(resource, payload, &result)
	if err != nil {
		return err
	}

	stat := ozon.Statistic{}
	stat.UUID = result.UUID

	c.storage.SetStatistic(stat)
	return nil
}
