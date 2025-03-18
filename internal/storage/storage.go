// Локальное хранилище
// Используется для промежуточного хранения данных
// Необходимо, так как генерация отчетов выполняется в три этапа
// 1. Запрос на формирование
// 2. Проверка готовности
// 3. Получение результатов

package storage

import (
	"fmt"
	"os"
	"ozonadv/internal/ozon"
)

type Storage struct {
	rootDir        string
	statisticsFile string
	statistics     map[string]ozon.Statistic
}

func New() *Storage {
	fmt.Println("Инициализация локального хранилища")

	root := os.TempDir() + "/ozonadv"
	s := Storage{
		rootDir:        root,
		statisticsFile: root + "/statistics.json",
		statistics:     make(map[string]ozon.Statistic),
	}

	initDir(s.rootDir)
	initFile(s.statisticsFile)
	readJsonFile(s.statisticsFile, &s.statistics, "{}")

	fmt.Println("Директория локального хранилища", s.rootDir)

	return &s
}

func (s *Storage) SetStatistic(item ozon.Statistic) {
	s.statistics[item.UUID] = item
}

func (s *Storage) GetStatistic(uuid string) (ozon.Statistic, bool) {
	item, ok := s.statistics[uuid]
	return item, ok
}

func (s *Storage) GetAllStatistic() []ozon.Statistic {
	result := make([]ozon.Statistic, 0, len(s.statistics))
	for _, item := range s.statistics {
		result = append(result, item)
	}

	return result
}

func (s *Storage) RemoveStatistic(uuid string) {
	delete(s.statistics, uuid)
}

func (s *Storage) StatisticsSize() int {
	return len(s.statistics)
}

// Сохранить состояние хранилища
func (s *Storage) SaveState() {
	fmt.Println("Сохранение локального хранилища")

	writeJsonFile(s.statisticsFile, s.statistics)
}
