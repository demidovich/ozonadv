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

	fmt.Println("Найдено кампаний:", s.storage.CampaignRequestsSize())
	if console.Ask("Удалить?") == false {
		return nil
	}

	s.storage.Reset()
	fmt.Println("Кампании удалены")

	return nil
}
