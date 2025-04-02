package object_stat

import (
	"fmt"
	"ozonadv/internal/storage"
)

type statInfoUsecase struct {
	storage *storage.Storage
}

func (s *statInfoUsecase) Handle() error {
	if s.storage.ObjectStatCampaigns().Size() == 0 {
		fmt.Println("Данные отсутствуют")
		fmt.Println("Формирование отчета не запускалось, либо была вызвана команда object_stat:reset")
		return nil
	}

	printReportInfo(s.storage)

	fmt.Println("")
	fmt.Println("ozonadv object_stat:continue Продолжить формирования отчета")
	fmt.Println("ozonadv object_stat:reset    Удалить незавершенное формирование отчета")

	return nil
}
