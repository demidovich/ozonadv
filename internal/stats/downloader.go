package stats

import (
	"errors"
	"fmt"
	"io"
	"os"
	"ozonadv/internal/ozon"
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
	out     io.Writer
	stat    *Stat
	storage storage
	ozon    *ozon.Ozon
}

func newDownloader(out io.Writer, stat *Stat, o *ozon.Ozon, s storage) *downloader {
	return &downloader{
		out:     out,
		stat:    stat,
		ozon:    o,
		storage: s,
	}
}

func (d *downloader) Start() {
	items := make(chan statItem)

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

func (d *downloader) createStatRequestsStage(in <-chan statItem) <-chan ozon.StatRequest {
	out := make(chan ozon.StatRequest)

	go func() {
		defer close(out)
		for item := range in {
			var statRequest ozon.StatRequest
			var err error
			var campaign ozon.Campaign = item.Campaign

			// Если ранее отчет формировался и закончился с ошибками либо был остановлен
			// Для некоторых кампаний уже были отправлены запросы для формирования статистики
			// Здесь мы проверяем есть ли они и используем их, если они есть

			if item.Request.UUID != "" {
				statRequest, err = d.ozon.StatRequests().Retrieve(item.Request.UUID)
				if err != nil {
					d.logCampaign(campaign, "уже есть сформированный запрос #", item.Request.UUID)
					d.logCampaign(campaign, "ошибка получения запроса: ", err)
					continue
				}
			} else {
				statRequest = d.createStatRequest(campaign)
				item.Request.UUID = statRequest.UUID
			}

			out <- statRequest
		}
	}()

	return out
}

func (d *downloader) readyStatRequestsStage(in <-chan ozon.StatRequest) <-chan ozon.StatRequest {
	out := make(chan ozon.StatRequest)

	go func() {
		defer close(out)
		for r := range in {
			var statRequest ozon.StatRequest
			var err error

			// Если ранее отчет формировался и закончился с ошибками либо был остановлен
			// Для некоторых кампаний уже были отправлены запросы для формирования статистики
			// Здесь мы проверяем в каком они состоянии и используем их, если с ними все ок

			if r.IsReadyToDownload() {
				statRequest = r
			} else {
				statRequest, err = d.readyStatRequest(r)
				if err != nil {
					d.logStatRequest(r, err)
					continue
				}
			}

			item, ok := d.stat.ItemByRequestUUID(statRequest.UUID)
			if !ok {
				d.logStatRequest(r, "не найдена кампания для запроса статистики")
				continue
			}

			item.Request.Link = statRequest.Link
			out <- statRequest
		}
	}()

	return out
}

func (d *downloader) downloadStatsStage(in <-chan ozon.StatRequest) <-chan bool {
	complete := make(chan bool)

	go func() {
		defer close(complete)
		for statRequest := range in {
			item, ok := d.stat.ItemByRequestUUID(statRequest.UUID)
			if !ok {
				d.logStatRequest(statRequest, "пропуск: не найдена кампания!!!")
				continue
			}

			if item.Request.File != "" {
				d.logStatRequest(statRequest, "статистика скачана ранее")
				continue
			}

			filename, err := d.downloadStat(statRequest)
			if err == nil {
				item.Request.File = filename
			} else {
				d.logStatRequest(statRequest, err)
			}
		}
		complete <- true
	}()

	return complete
}

func (d *downloader) createStatRequest(campaign ozon.Campaign) ozon.StatRequest {
	options := ozon.CreateStatRequestOptions{
		CampaignId: campaign.ID,
		DateFrom:   d.stat.Options.DateFrom,
		DateTo:     d.stat.Options.DateTo,
		GroupBy:    d.stat.Options.GroupBy,
	}

	for attempt := 1; attempt <= createStatRequestAttempts; attempt++ {
		d.logCampaign(campaign, "создание запроса отчета: попытка ", attempt)

		var err error
		var req ozon.StatRequest

		if d.stat.Options.Type == "TOTAL" {
			req, err = d.ozon.StatRequests().CreateTotal(campaign, options)
		} else {
			req, err = d.ozon.StatRequests().CreateObject(campaign, options)
		}

		if err == nil {
			d.logCampaign(campaign, "создан запрос ", req.UUID)
			return req
		}

		d.logCampaign(campaign, err)

		waitTime := createStatRequestWaitTime
		waitTime = waitTime + waitTime*time.Duration(attempt-1)
		d.logCampaign(campaign, "создание запроса отчета: ждем ", waitTime.String())
		time.Sleep(waitTime)
	}

	d.logString("Превышено количество количество попыток создания запроса.\n")
	d.logString("Возможно существует тяжелый несформированый запрос.\n")
	d.logString("Пока Озон не закончит его формирование создать новый не получится.\n")
	d.logString("Выдержите паузу и запустите ozonadv stat:continue.\n")
	d.logString("")
	os.Exit(1)

	return ozon.StatRequest{}
}

func (d *downloader) readyStatRequest(statRequest ozon.StatRequest) (ozon.StatRequest, error) {
	for attempt := 1; attempt <= readyStatRequestAttempts; attempt++ {
		waitTime := readyStatRequestWaitTime

		d.logStatRequest(statRequest, "ожидание готовности: ждем ", waitTime.String())
		time.Sleep(waitTime)

		d.logStatRequest(statRequest, "ожидание готовности: попытка ", attempt)

		// Сначала ждем, так как после создания запрос будет собираться

		req, err := d.ozon.StatRequests().Retrieve(statRequest.UUID)
		if err != nil {
			d.logStatRequest(statRequest, "ожидание готовности: ", err)
			continue
		}

		if !req.IsReadyToDownload() {
			d.logStatRequest(statRequest, "ожидание готовности: состояние ", req.State)
			continue
		}

		d.logStatRequest(statRequest, "готов к скачиванию")

		return req, nil
	}

	return ozon.StatRequest{}, errors.New("превышено количество попыток")
}

func (d *downloader) downloadStat(statRequest ozon.StatRequest) (string, error) {
	filename := "object-stat-" + statRequest.CampaignId() + ".csv"
	for attempt := 1; attempt <= downloadStatAttempts; attempt++ {
		d.logStatRequest(statRequest, "скачивание статистики: попытка ", attempt)

		data, err := d.ozon.StatRequests().Download(statRequest)
		if err != nil {
			d.logStatRequest(statRequest, "скачивание статистики: ", err)
			time.Sleep(downloadStatWaitTime)
			continue
		}

		err = d.storage.Downloads().Write(filename, data)
		if err != nil {
			d.logStatRequest(statRequest, "скачивание статистики: ", err)
			time.Sleep(downloadStatWaitTime)
			continue
		}

		d.logStatRequest(statRequest, "скачан файл: ", filename)
		return filename, nil
	}

	return "", errors.New("превышено количество попыток скачивания")
}

func (d *downloader) logCampaign(c ozon.Campaign, msg ...any) {
	d.logString(
		fmt.Sprintf("[%s] %s\n", c.ID, fmt.Sprint(msg...)),
	)
}

func (d *downloader) logStatRequest(r ozon.StatRequest, msg ...any) {
	d.logString(
		fmt.Sprintf("[%s] %s\n", r.CampaignId(), fmt.Sprint(msg...)),
	)
}

func (d *downloader) logString(msg string) {
	d.out.Write([]byte(msg))
}
