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
	Campaigns            *campaigns
	requestOptionsFile   string
	requestOptions       *RequestOptions
	processedRequestFile string
	processedRequest     *ozon.StatisticRequest
	downloadsDir         string
	Downloads            Downloads
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
		downloadsDir:         root + "/downloads",
	}

	utils.DirInitOrFail(s.rootDir)
	utils.FileInitOrFail(s.campanignsFile)
	utils.FileInitOrFail(s.requestOptionsFile)
	utils.FileInitOrFail(s.processedRequestFile)
	utils.JsonFileReadOrFail(s.requestOptionsFile, &s.requestOptions, "{}")
	utils.JsonFileReadOrFail(s.processedRequestFile, &s.processedRequest, "{}")
	utils.DirInitOrFail(s.downloadsDir)

	s.Campaigns = NewCampaigns(s.campanignsFile)
	s.Downloads = NewDownloads(s.downloadsDir)

	return &s
}

func (s *Storage) RootDir() string {
	return s.rootDir
}

func (s *Storage) SetRequestOptions(options RequestOptions) {
	s.requestOptions = &options
}

func (s *Storage) RequestOptions() *RequestOptions {
	return s.requestOptions
}

func (s *Storage) SetProcessedRequest(request *ozon.StatisticRequest) {
	s.processedRequest = request
}

func (s *Storage) ProcessedRequest() *ozon.StatisticRequest {
	return s.processedRequest
}

// Reset all storage data
func (s *Storage) Reset() error {
	s.Campaigns.RemoveAll()
	s.requestOptions = nil
	s.processedRequest = nil

	return s.Downloads.RemoveAll()
}

// Сохранить состояние хранилища
func (s *Storage) SaveState() {
	fmt.Println("")
	fmt.Println("Сохранение локального хранилища")

	err := utils.JsonFileWrite(s.campanignsFile, s.Campaigns.All())
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
