package stats

import (
	"errors"
	"fmt"
	"os"
	"ozonadv/internal/infra/ozon"
	"ozonadv/internal/models"
	"sync"
	"time"
)

const (
	createStatRequestAttempts = 5
	createStatRequestWaitTime = 30 * time.Second
	readyStatRequestAttempts  = 10
	readyStatRequestWaitTime  = 30 * time.Second
	downloadStatAttempts      = 5
	downloadStatWaitTime      = 10 * time.Second
)

type downloader struct {
	stat    *models.Stat
	storage Storage
	ozon    *ozon.Ozon
	debug   Debug
	saveMu  *sync.Mutex
}

func newDownloader(stat *models.Stat, ozon *ozon.Ozon, storage Storage, debug Debug) *downloader {
	return &downloader{
		stat:    stat,
		ozon:    ozon,
		storage: storage,
		debug:   debug,
		saveMu:  &sync.Mutex{},
	}
}

func (d *downloader) Start() {
	d.stat.RunnedAt = time.Now().String()

	d.storage.Add(d.stat)
	defer d.storage.Add(d.stat)

	items := make(chan models.StatItem)

	go func() {
		defer close(items)
		for _, item := range d.stat.Items {
			items <- item
		}
	}()

	createdRequests := d.createStatRequestsStage(items)
	readyRequests := d.readyStatRequestsStage(createdRequests)
	<-d.downloadStatsStage(readyRequests)
}

func (d *downloader) createStatRequestsStage(in <-chan models.StatItem) <-chan models.StatRequest {
	out := make(chan models.StatRequest)

	go func() {
		defer close(out)
		for item := range in {
			var statRequest models.StatRequest
			var err error
			var campaign models.Campaign = item.Campaign

			// Если ранее отчет формировался и закончился с ошибками либо был остановлен
			// Для некоторых кампаний уже были отправлены запросы для формирования статистики
			// Здесь мы проверяем есть ли они и используем их, если они есть

			if item.Request.UUID != "" {
				statRequest, err = d.ozon.StatRequests().Retrieve(item.Request.UUID)
				if err != nil {
					d.debugCampaign(campaign, "уже есть сформированный запрос #", item.Request.UUID)
					d.debugCampaign(campaign, "ошибка получения запроса: ", err)
					continue
				}
			} else {
				statRequest = d.createStatRequest(campaign)
				item.Request.UUID = statRequest.UUID
			}

			d.saveStatRequest(statRequest)
			out <- statRequest
		}
	}()

	return out
}

func (d *downloader) readyStatRequestsStage(in <-chan models.StatRequest) <-chan models.StatRequest {
	out := make(chan models.StatRequest)

	go func() {
		defer close(out)
		for r := range in {
			var statRequest models.StatRequest
			var err error

			// Если ранее отчет формировался и закончился с ошибками либо был остановлен
			// Для некоторых кампаний уже были отправлены запросы для формирования статистики
			// Здесь мы проверяем в каком они состоянии и используем их, если с ними все ок

			if r.IsReadyToDownload() {
				statRequest = r
			} else {
				statRequest, err = d.readyStatRequest(r)
				if err != nil {
					d.debugStatRequest(r, err)
					continue
				}
			}

			item, ok := d.stat.ItemByRequestUUID(statRequest.UUID)
			if !ok {
				d.debugStatRequest(r, "не найдена кампания для запроса статистики")
				continue
			}

			item.Request.Link = statRequest.Link
			d.saveStatRequest(statRequest)
			out <- statRequest
		}
	}()

	return out
}

func (d *downloader) downloadStatsStage(in <-chan models.StatRequest) <-chan bool {
	complete := make(chan bool)

	go func() {
		defer close(complete)
		for statRequest := range in {
			item, ok := d.stat.ItemByRequestUUID(statRequest.UUID)
			if !ok {
				d.debugStatRequest(statRequest, "пропуск: не найдена кампания!!!")
				continue
			}

			if item.Request.File != "" {
				d.debugStatRequest(statRequest, "статистика скачана ранее")
				continue
			}

			filename, err := d.downloadStat(statRequest)
			if err == nil {
				item.Request.File = filename
				d.saveStatRequestFile(statRequest, filename)
			} else {
				d.debugStatRequest(statRequest, err)
			}
		}
		complete <- true
	}()

	return complete
}

func (d *downloader) createStatRequest(campaign models.Campaign) models.StatRequest {
	startedAt := time.Now()

	options := ozon.CreateStatRequestOptions{
		CampaignId: campaign.ID,
		DateFrom:   d.stat.Options.DateFrom,
		DateTo:     d.stat.Options.DateTo,
		GroupBy:    d.stat.Options.GroupBy,
	}

	for attempt := 1; attempt <= createStatRequestAttempts; attempt++ {
		d.debugCampaignTime(
			campaign,
			time.Since(startedAt),
			"создание запроса отчета: попытка ",
			attempt,
		)

		var err error
		var req models.StatRequest

		if d.stat.Options.Type == "TOTAL" {
			req, err = d.ozon.StatRequests().CreateTotal(campaign, options)
		} else {
			req, err = d.ozon.StatRequests().CreateObject(campaign, options)
		}

		if err == nil {
			d.debugCampaign(campaign, "создан запрос ", req.UUID)
			return req
		}

		d.debugCampaign(campaign, err)

		waitTime := createStatRequestWaitTime
		waitTime = waitTime + waitTime*time.Duration(attempt-1)
		d.debugCampaign(campaign, "создание запроса отчета: ждем ", waitTime.String())
		time.Sleep(waitTime)
	}

	d.debug.Println("Превышено количество количество попыток создания запроса.")
	d.debug.Println("Возможно существует тяжелый несформированый запрос.")
	d.debug.Println("Пока Озон не закончит его формирование создать новый не получится.")
	d.debug.Println("Выдержите паузу и повторно запустите формирование этого отчета.")
	d.debug.Println("")
	os.Exit(1)

	return models.StatRequest{}
}

func (d *downloader) readyStatRequest(statRequest models.StatRequest) (models.StatRequest, error) {
	startedAt := time.Now()

	for attempt := 1; attempt <= readyStatRequestAttempts; attempt++ {
		waitTime := readyStatRequestWaitTime

		d.debugStatRequest(statRequest, "ожидание готовности: ждем ", waitTime.String())
		time.Sleep(waitTime)

		d.debugStatRequestTime(
			statRequest,
			time.Since(startedAt),
			"ожидание готовности: попытка ",
			attempt,
		)

		// Сначала ждем, так как после создания запрос будет собираться

		req, err := d.ozon.StatRequests().Retrieve(statRequest.UUID)
		if err != nil {
			d.debugStatRequest(statRequest, "ожидание готовности: ", err)
			continue
		}

		if !req.IsReadyToDownload() {
			d.debugStatRequest(statRequest, "ожидание готовности: состояние ", req.State)
			continue
		}

		d.debugStatRequestTime(
			statRequest,
			time.Since(startedAt),
			"готов к скачиванию, время ",
		)

		return req, nil
	}

	return models.StatRequest{}, errors.New("превышено количество попыток")
}

func (d *downloader) downloadStat(statRequest models.StatRequest) (string, error) {
	filename := "object-stat-" + statRequest.CampaignId() + ".csv"
	for attempt := 1; attempt <= downloadStatAttempts; attempt++ {
		d.debugStatRequest(statRequest, "скачивание статистики: попытка ", attempt)

		data, err := d.ozon.StatRequests().Download(statRequest)
		if err != nil {
			d.debugStatRequest(statRequest, "скачивание статистики: ", err)
			time.Sleep(downloadStatWaitTime)
			continue
		}

		d.storage.SaveDownloadedFile(d.stat, filename, data)
		d.debugStatRequest(statRequest, "скачан файл: ", filename)

		return filename, nil
	}

	return "", errors.New("превышено количество попыток скачивания")
}

func (d *downloader) saveStatRequest(statRequest models.StatRequest) {
	d.saveMu.Lock()
	defer d.saveMu.Unlock()

	campaignID := statRequest.CampaignId()
	for i, item := range d.stat.Items {
		if item.Campaign.ID != campaignID {
			continue
		}
		d.stat.Items[i].Request.UUID = statRequest.UUID
		d.stat.Items[i].Request.Link = statRequest.Link
	}

	d.storage.Add(d.stat)
}

func (d *downloader) saveStatRequestFile(statRequest models.StatRequest, file string) {
	d.saveMu.Lock()
	defer d.saveMu.Unlock()

	campaignID := statRequest.CampaignId()
	for i, item := range d.stat.Items {
		if item.Campaign.ID != campaignID {
			continue
		}
		d.stat.Items[i].Request.File = file
	}

	d.storage.Add(d.stat)
}

func (d *downloader) debugCampaign(c models.Campaign, msg ...any) {
	d.debug.Printf("[%s] %s\n", c.ID, fmt.Sprint(msg...))
}

func (d *downloader) debugCampaignTime(c models.Campaign, t time.Duration, msg ...any) {
	tString, _ := time.Parse(time.TimeOnly, t.String())
	d.debug.Printf("[%s] %s, время %s\n", c.ID, fmt.Sprint(msg...), tString)
}

func (d *downloader) debugStatRequest(r models.StatRequest, msg ...any) {
	d.debug.Printf("[%s] %s\n", r.CampaignId(), fmt.Sprint(msg...))
}

func (d *downloader) debugStatRequestTime(r models.StatRequest, t time.Duration, msg ...any) {
	tString, _ := time.Parse(time.TimeOnly, t.String())
	d.debug.Printf("[%s] %s, время %s\n", r.CampaignId(), fmt.Sprint(msg...), tString)
}
