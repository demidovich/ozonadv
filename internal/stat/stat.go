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

	"github.com/go-playground/validator/v10"
)

type HandleOptions struct {
	FromDate string `validate:"required,datetime=2006-01-02"`
	ToDate   string `validator:"required,datetime=2006-01-02"`
}

func (f *HandleOptions) validate() error {
	validate := validator.New()

	err := validate.Struct(f)
	if err == nil {
		return nil
	}

	errs := err.(validator.ValidationErrors)
	return fmt.Errorf("%s", errs)
}

func Handle(ozonClient *ozon.Client, options HandleOptions) error {
	if err := options.validate(); err != nil {
		return err
	}

	fmt.Println(1111)
	os.Exit(1)

	campaigns, err := ozonClient.Campaigns()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("")
	fmt.Printf("Найдено компаний: %d\n", len(campaigns))
	fmt.Println("")

	for _, campaign := range campaigns {
		fmt.Printf("#%s, %s, %s\n", campaign.ID, campaign.State, campaign.Title)

		if campaign.NotRunning() {
			fmt.Println("Пропуск. Кампания не была запущена.")
			fmt.Println("")
			continue
		}

		// objects, err := ozonClient.CampaignObjects(campaign.ID)
		// if err != nil {
		// 	fmt.Println(err)
		// 	continue
		// }

		// if len(objects) == 0 {
		// 	fmt.Println("Пропуск. Отсутствуют рекламные объекты.")
		// 	fmt.Println("")
		// 	continue
		// }
	}

	return nil
}
