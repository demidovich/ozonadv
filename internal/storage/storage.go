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
	rootDir                  string
	statOptionsFile          string
	statOptions              *StatOptions
	statCampanignsFile       string
	statCampaigns            *campaigns
	objectStatOptionsFile    string
	objectStatOptions        *ObjectStatOptions
	objectStatCampanignsFile string
	objectStatCampaigns      *campaigns
	downloadsDir             string
	downloads                downloads
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

type ObjectStatOptions struct {
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
		rootDir:                  root,
		statOptionsFile:          root + "/stat-options.json",
		statCampanignsFile:       root + "/stat-campaigns.json",
		objectStatOptionsFile:    root + "/object-stat-options.json",
		objectStatCampanignsFile: root + "/object-stat-campaigns.json",
		downloadsDir:             root + "/downloads",
	}

	utils.DirInitOrFail(s.rootDir)
	utils.FileInitOrFail(s.statOptionsFile)
	utils.FileInitOrFail(s.statCampanignsFile)
	utils.FileInitOrFail(s.objectStatOptionsFile)
	utils.FileInitOrFail(s.objectStatCampanignsFile)
	utils.DirInitOrFail(s.downloadsDir)
	utils.JsonFileReadOrFail(s.statOptionsFile, &s.statOptions, "{}")
	utils.JsonFileReadOrFail(s.objectStatOptionsFile, &s.objectStatOptions, "{}")

	s.statCampaigns = NewCampaigns(s.statCampanignsFile)
	s.objectStatCampaigns = NewCampaigns(s.objectStatCampanignsFile)
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

func (s *Storage) StatCampaigns() *campaigns {
	return s.statCampaigns
}

func (s *Storage) ObjectStatCampaigns() *campaigns {
	return s.statCampaigns
}

func (s *Storage) Downloads() downloads {
	return s.downloads
}

// Reset all storage data
func (s *Storage) Reset() error {
	s.statOptions = nil
	s.statCampaigns.RemoveAll()
	s.downloads.RemoveAll()

	return s.downloads.RemoveAll()
}

// Сохранить состояние хранилища
func (s *Storage) SaveState() {
	fmt.Println("")
	fmt.Println("[shutdown] сохранение локального хранилища")

	if err := utils.JsonFileWrite(s.statOptionsFile, s.statOptions); err != nil {
		log.Fatal(err)
	}

	if err := utils.JsonFileWrite(s.statCampanignsFile, s.statCampaigns.Data()); err != nil {
		log.Fatal(err)
	}

	if err := utils.JsonFileWrite(s.objectStatOptionsFile, s.objectStatOptions); err != nil {
		log.Fatal(err)
	}

	if err := utils.JsonFileWrite(s.objectStatCampanignsFile, s.objectStatCampaigns.Data()); err != nil {
		log.Fatal(err)
	}
}
