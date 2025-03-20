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
	"ozonadv/internal/ozon"
	"ozonadv/pkg/utils"
)

type Storage struct {
	rootDir              string
	campanignsFile       string
	campaigns            map[string]ozon.Campaign
	requestOptionsFile   string
	requestOptions       *RequestOptions
	processedRequestFile string
	processedRequest     *ozon.StatisticRequest
}

type RequestOptions struct {
	DateFrom   string `json:"dateFrom"`
	DateTo     string `json:"dateTo"`
	ExportFile string `json:"exportFile"`
	GroupBy    string `json:"groupBy"`
}

func New() *Storage {
	root := os.TempDir() + "/ozonadv"
	s := Storage{
		rootDir:              root,
		campanignsFile:       root + "/campaigns.json",
		requestOptionsFile:   root + "/request-options.json",
		processedRequestFile: root + "/processed-request.json",
		campaigns:            make(map[string]ozon.Campaign),
	}

	utils.DirInitOrFail(s.rootDir)
	utils.FileInitOrFail(s.campanignsFile)
	utils.FileInitOrFail(s.requestOptionsFile)
	utils.FileInitOrFail(s.processedRequestFile)
	utils.JsonFileReadOrFail(s.campanignsFile, &s.campaigns, "{}")
	utils.JsonFileReadOrFail(s.requestOptionsFile, &s.requestOptions, "{}")
	utils.JsonFileReadOrFail(s.processedRequestFile, &s.processedRequest, "{}")

	return &s
}

func (s *Storage) RootDir() string {
	return s.rootDir
}

func (s *Storage) AddCampaignRequest(item ozon.Campaign) {
	s.campaigns[item.ID] = item
}

func (s *Storage) HasCampaignRequest(id string) bool {
	_, ok := s.campaigns[id]
	return ok
}

func (s *Storage) CampaignRequests() []ozon.Campaign {
	result := make([]ozon.Campaign, 0, len(s.campaigns))
	for _, item := range s.campaigns {
		result = append(result, item)
	}

	return result
}

func (s *Storage) CampaignRequestsSize() int {
	return len(s.campaigns)
}

func (s *Storage) RemoveCampaignRequest(id string) {
	delete(s.campaigns, id)
}

func (s *Storage) SetRequestOptions(options RequestOptions) {
	s.requestOptions = &options
}

func (s *Storage) RequestOptions() *RequestOptions {
	return s.requestOptions
}

func (s *Storage) ProcessedRequest() *ozon.StatisticRequest {
	return s.processedRequest
}

// Reset all storage data
func (s *Storage) Reset() {
	for k := range s.campaigns {
		delete(s.campaigns, k)
	}
	s.requestOptions = nil
	s.processedRequest = nil
}

// Сохранить состояние хранилища
func (s *Storage) SaveState() {
	fmt.Println("")
	fmt.Println("Сохранение локального хранилища")

	err := utils.JsonFileWrite(s.campanignsFile, s.campaigns)
	if err != nil {
		log.Fatal(err)
	}

	err = utils.JsonFileWrite(s.requestOptionsFile, s.requestOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = utils.JsonFileWrite(s.processedRequestFile, s.processedRequest)
	if err != nil {
		log.Fatal(err)
	}
}
