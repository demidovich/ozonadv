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
	"fmt"
	"ozonadv/internal/ozon"
	"ozonadv/internal/storage"
	"time"

	"github.com/go-playground/validator/v10"
)

type CreateOptions struct {
	FromDate   string `validate:"required,datetime=2006-01-02"`
	ToDate     string `validate:"required,datetime=2006-01-02"`
	ExportFile string `validate:"required,filepath"`
}

func (c *CreateOptions) validate() error {
	validate := validator.New()

	err := validate.Struct(c)
	if err == nil {
		return nil
	}

	errs := err.(validator.ValidationErrors)
	if errs != nil {
		return fmt.Errorf("%s", errs)
	}

	return nil
}

type createUsecase struct {
	storage     *storage.Storage
	ozonClient  *ozon.Client
	pullUsecase pullUsecase
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

	if len(campaigns) == 0 {
		fmt.Println("Кампании не найдены.")
		return nil
	}

	fmt.Printf("Найдено кампаний: %d\n", len(campaigns))
	c.createStatisticRequests(campaigns, options)

	if c.storage.StatisticsSize() == 0 {
		fmt.Println("Ни одного запроса не сформировано.")
		return nil
	}

	// Прогрессирующий таймаут для повторных опросов
	// Необходимо на случай, если выгрузка будет объемной или залипнет

	attempt := 0
	retries := 3
	timeout := 3 * time.Second
	time.Sleep(timeout)

	for {
		if attempt == retries {
			return errors.New("Превышен лимит повторных попыток загрузки отчетов.")
		}

		if c.storage.StatisticsSize() == 0 {
			break
		}

		attempt++
		timeout *= 2
		time.Sleep(timeout)

		pullOptions := PullOptions{ExportFile: options.ExportFile}
		c.pullUsecase.Handle(pullOptions)
	}

	fmt.Println("Все отчеты загружены.")
	return nil
}

func (c *createUsecase) createStatisticRequests(campaigns []ozon.Campaign, options CreateOptions) {
	fmt.Println("Формирование запросов на генерацию статистики.")
	fmt.Println("")

	for _, campaign := range campaigns {
		fmt.Printf("#%s, %s, %s\n", campaign.ID, campaign.State, campaign.Title)

		if campaign.NotRunning() {
			fmt.Println("Пропуск. Кампания не была запущена.")
			fmt.Println("")
			continue
		}

		err := c.requestCampaignStatistic(campaign, options)
		if err != nil {
			fmt.Println("Ошибка запроса формирования отчета:", err)
			continue
		}

		fmt.Println("Отправлен запрос формирования отчета.")
		fmt.Println("")
	}
}

func (c *createUsecase) requestCampaignStatistic(campaign ozon.Campaign, options CreateOptions) error {
	var resource string
	var payload map[string]any

	// Запрос UI
	//
	// if campaign.AdvObjectType == "VIDEO_BANNER" {
	// 	resource = "/adv-api/external/api/statistics/video"
	// 	payload = map[string]any{
	// 		"campaigns": []string{campaign.ID},
	// 		"dateFrom":  options.FromDate,
	// 		"dateTo":    options.ToDate,
	// 		"groupBy":   "DATE",
	// 	}
	// } else {
	// 	resource = "/adv-api/external/api/statistics"
	// 	payload = map[string]any{
	// 		"campaignId": campaign.ID,
	// 		"dateFrom":   options.FromDate,
	// 		"dateTo":     options.ToDate,
	// 		"groupBy":    "DATE",
	// 	}
	// }

	if campaign.AdvObjectType == "VIDEO_BANNER" {
		resource = "/client/statistics/video"
	} else {
		resource = "/client/statistics"
	}

	payload = map[string]any{
		"campaigns": []string{campaign.ID},
		"dateFrom":  options.FromDate,
		"dateTo":    options.ToDate,
		"groupBy":   "DATE",
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
