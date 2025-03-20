package stat

import (
	"fmt"
	"ozonadv/internal/ozon"
	"ozonadv/internal/storage"
)

type statInfoUsecase struct {
	storage    *storage.Storage
	ozonClient *ozon.Client
}

func (s *statInfoUsecase) Handle() error {
	if s.storage.CampaignRequestsSize() > 0 {
		fmt.Println("Есть незавершенное формирование отчета")
		fmt.Println("Кампаний для обработки:", s.storage.CampaignRequestsSize())
		fmt.Println("")
		fmt.Println("ozonadv stat:continue Продолжить формирования отчета")
		fmt.Println("ozonadv stat:reset    Удалить незавершенное формирование отчета")
	} else {
		fmt.Println("Формирование отчета завершено")
	}

	return nil
}
