package stat

import (
	"fmt"
	"ozonadv/internal/storage"
	"ozonadv/pkg/console"
)

type statResetUsecase struct {
	storage *storage.Storage
}

func (s *statResetUsecase) Handle() error {
	if s.storage.Campaigns.Size() == 0 {
		fmt.Println("Нет необработанных кампаний")
		return nil
	}

	fmt.Println("Найдено необработанных кампаний:", s.storage.Campaigns.Size())
	fmt.Println("")
	if console.Ask("Удалить?") == false {
		return nil
	}
	fmt.Println("")
	fmt.Println("Кампании удалены")

	return s.storage.Reset()
}
