// Локальное хранилище
// Используется для промежуточного хранения данных
// Необходимо, так как генерация отчетов выполняется в три этапа
// 1. Запрос на формирование
// 2. Проверка готовности
// 3. Получение результатов

package storage

import (
	"fmt"
	"log"
	"os"
	"ozonadv/pkg/utils"
)

type Storage struct {
	rootDir         string
	statOptionsFile string
	statOptions     *StatOptions
	campanignsFile  string
	campaigns       *campaigns
	downloadsDir    string
	downloads       downloads
}

// По сути это stat.StatOptions
// Но импортировать ее нельзя из-за циклической зависомости
type StatOptions struct {
	DateFrom         string `json:"dateFrom"`
	DateTo           string `json:"dateTo"`
	GroupBy          string `json:"groupBy"`
	CreatedAt        string `json:"createdAt"`
	StartedAt        string `json:"startedAt"`
	ApiRequestsCount int    `json:"apiRequestsCount"`
}

func New() *Storage {
	root := os.TempDir() + "/ozonadv"
	s := Storage{
		rootDir:         root,
		statOptionsFile: root + "/stat-options.json",
		campanignsFile:  root + "/campaigns.json",
		downloadsDir:    root + "/downloads",
	}

	utils.DirInitOrFail(s.rootDir)
	utils.FileInitOrFail(s.statOptionsFile)
	utils.FileInitOrFail(s.campanignsFile)
	utils.DirInitOrFail(s.downloadsDir)
	utils.JsonFileReadOrFail(s.statOptionsFile, &s.statOptions, "{}")

	s.campaigns = NewCampaigns(s.campanignsFile)
	s.downloads = NewDownloads(s.downloadsDir)

	return &s
}

func (s *Storage) RootDir() string {
	return s.rootDir
}

func (s *Storage) SetStatOptions(options StatOptions) {
	s.statOptions = &options
}

func (s *Storage) StatOptions() *StatOptions {
	return s.statOptions
}

func (s *Storage) Campaigns() *campaigns {
	return s.campaigns
}

func (s *Storage) Downloads() downloads {
	return s.downloads
}

// Reset all storage data
func (s *Storage) Reset() error {
	s.campaigns.RemoveAll()
	s.downloads.RemoveAll()
	s.statOptions = nil

	return s.downloads.RemoveAll()
}

// Сохранить состояние хранилища
func (s *Storage) SaveState() {
	fmt.Println("")
	fmt.Println("[shutdown] сохранение локального хранилища")

	err := utils.JsonFileWrite(s.campanignsFile, s.campaigns.Data())
	if err != nil {
		log.Fatal(err)
	}

	err = utils.JsonFileWrite(s.statOptionsFile, s.statOptions)
	if err != nil {
		log.Fatal(err)
	}
}
