package stat

import (
	"fmt"
	"ozonadv/internal/storage"
)

type statInfoUsecase struct {
	storage *storage.Storage
}

func (s *statInfoUsecase) Handle() error {
	if s.storage.StatCampaigns().Size() > 0 {
		fmt.Println("Есть незавершенное формирование отчета")
		printCampaignsTable(s.storage.StatCampaigns().All())
		fmt.Println("")
		fmt.Println("ozonadv stat:continue Продолжить формирования отчета")
		fmt.Println("ozonadv stat:reset    Удалить незавершенное формирование отчета")
	} else {
		fmt.Println("Данные отсутствуют")
		fmt.Println("Формирование отчета не запускалось, либо была вызвана команда stat:reset")
	}

	return nil
}
