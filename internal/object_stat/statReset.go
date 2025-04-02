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
	if s.storage.ObjectStatCampaigns().Size() == 0 {
		printOptionsTable(*s.storage.ObjectStatOptions())
	} else {
		printReportInfo(s.storage)
	}

	fmt.Println("")
	if console.Ask("Удалить параметры формирования отчета?") == false {
		return nil
	}

	err := s.storage.ObjectStatResetAll()
	if err != nil {
		return err
	}

	fmt.Println("")
	fmt.Println("Параметры формирования отчета удалены")

	return nil
}
