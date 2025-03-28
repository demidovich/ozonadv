package object_stat

import (
	"fmt"
	"ozonadv/internal/storage"
	"ozonadv/pkg/console"
)

type statResetUsecase struct {
	storage *storage.Storage
}

func (s *statResetUsecase) Handle() error {
	if s.storage.ObjectStatCampaigns().Size() > 0 {
		fmt.Println("Найдено кампаний:", s.storage.ObjectStatCampaigns().Size())
		fmt.Println("")
	}

	if console.Ask("Удалить параметры формирования отчета?") == false {
		return nil
	}

	fmt.Println("")
	fmt.Println("Параметры формирования отчета удалены")

	return s.storage.ObjectStatReset()
}
