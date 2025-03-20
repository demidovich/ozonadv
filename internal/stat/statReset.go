package stat

import (
	"fmt"
	"ozonadv/internal/ozon"
	"ozonadv/internal/storage"
	"ozonadv/pkg/console"
)

type statResetUsecase struct {
	storage    *storage.Storage
	ozonClient *ozon.Client
}

func (s *statResetUsecase) Handle() error {
	if s.storage.CampaignRequestsSize() == 0 {
		fmt.Println("Нет необработанных кампаний")
		return nil
	}

	fmt.Println("Найдено необработанных кампаний:", s.storage.CampaignRequestsSize())
	fmt.Println("")
	if console.Ask("Удалить?") == false {
		return nil
	}

	s.storage.Reset()
	fmt.Println("")
	fmt.Println("Кампании удалены")

	return nil
}
