package stat

import (
	"fmt"
	"os"
	"ozonstat/internal/ozon"
)

type FetchOptions struct {
	Days uint
}

func Fetch(ozonClient *ozon.Client, options FetchOptions) error {
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

		objects, err := ozonClient.CampaignObjects(campaign.ID)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if len(objects) == 0 {
			fmt.Println("Объекты не найдены")
			fmt.Println("")
			continue
		}

		for object := range objects {

		}
	}

	return nil
}
