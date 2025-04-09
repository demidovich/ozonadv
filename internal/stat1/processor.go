package stat1

import (
	"errors"
	"fmt"
	"io"
	"os"
	"ozonadv/internal/ozon"
	"ozonadv/internal/storage"
	"time"
)

const (
	createAttempts   = 5
	createWaitTime   = 30 * time.Second
	readyAttempts    = 10
	readyWaitTime    = 30 * time.Second
	downloadAttempts = 5
	downloadWaitTime = 10 * time.Second
)

type processor struct {
	stat    *Stat
	storage *storage.Storage
	ozon    *ozon.Ozon
	out     io.Writer
}

func newProcessor(out io.Writer, stat *Stat, o *ozon.Ozon, s *storage.Storage) *processor {
	return &processor{
		stat:    stat,
		ozon:    o,
		storage: s,
		out:     out,
	}
}

func (p *processor) Start() {
	items := make(chan statItem)

	go func() {
		defer close(items)
		for _, item := range p.stat.Items {
			items <- item
		}
	}()

	createdRequests := p.createStatRequestsStage(items)
	readyRequests := p.readyStatRequestsStage(createdRequests)
	<-p.downloadStatsStage(readyRequests)
}

func (p *processor) createStatRequestsStage(in <-chan statItem) <-chan ozon.StatRequest {
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
				statRequest, err = p.ozon.StatRequests().Retrieve(item.Request.UUID)
				if err != nil {
					p.logCampaign(campaign, "уже есть сформированный запрос #", item.Request.UUID)
					p.logCampaign(campaign, "ошибка получения запроса: ", err)
					continue
				}
			} else {
				statRequest = p.createStatRequest(campaign)
				item.Request.UUID = statRequest.UUID
			}

			out <- statRequest
		}
	}()

	return out
}

func (p *processor) readyStatRequestsStage(in <-chan ozon.StatRequest) <-chan ozon.StatRequest {
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
				statRequest, err = p.readyStatRequest(r)
				if err != nil {
					p.logStatRequest(r, err)
					continue
				}
			}

			item, ok := p.stat.ItemByRequestUUID(statRequest.UUID)
			if !ok {
				p.logStatRequest(r, "не найдена кампания для запроса статистики")
				continue
			}

			item.Request.Link = statRequest.Link
			out <- statRequest
		}
	}()

	return out
}

func (p *processor) downloadStatsStage(in <-chan ozon.StatRequest) <-chan bool {
	complete := make(chan bool)

	go func() {
		defer close(complete)
		for statRequest := range in {
			item, ok := p.stat.ItemByRequestUUID(statRequest.UUID)
			if !ok {
				p.logStatRequest(statRequest, "пропуск: не найдена кампания!!!")
				continue
			}

			if item.Request.File != "" {
				p.logStatRequest(statRequest, "статистика скачана ранее")
				continue
			}

			filename, err := p.downloadStat(statRequest)
			if err == nil {
				item.Request.File = filename
			} else {
				p.logStatRequest(statRequest, err)
			}
		}
		complete <- true
	}()

	return complete
}

func (p *processor) createStatRequest(campaign ozon.Campaign) ozon.StatRequest {
	options := ozon.CreateStatRequestOptions{
		CampaignId: campaign.ID,
		DateFrom:   p.stat.Options.DateFrom,
		DateTo:     p.stat.Options.DateTo,
		GroupBy:    p.stat.Options.GroupBy,
	}

	for attempt := 1; attempt <= createAttempts; attempt++ {
		p.logCampaign(campaign, "создание запроса отчета: попытка ", attempt)

		var err error
		var req ozon.StatRequest

		if p.stat.Options.Type == "TOTAL" {
			req, err = p.ozon.StatRequests().CreateTotal(campaign, options)
		} else {
			req, err = p.ozon.StatRequests().CreateObject(campaign, options)
		}

		if err == nil {
			p.logCampaign(campaign, "создан запрос ", req.UUID)
			return req
		}

		p.logCampaign(campaign, err)

		waitTime := createWaitTime + createWaitTime*time.Duration(attempt-1)
		p.logCampaign(campaign, "создание запроса отчета: ждем ", waitTime.String())
		time.Sleep(waitTime)
	}

	p.logString("Превышено количество количество попыток создания запроса.\n")
	p.logString("Возможно существует тяжелый несформированый запрос.\n")
	p.logString("Пока Озон не закончит его формирование создать новый не получится.\n")
	p.logString("Выдержите паузу и запустите ozonadv stat:continue.\n")
	p.logString("")
	os.Exit(1)

	return ozon.StatRequest{}
}

func (p *processor) readyStatRequest(statRequest ozon.StatRequest) (ozon.StatRequest, error) {
	for attempt := 1; attempt <= readyAttempts; attempt++ {
		waitTime := readyWaitTime

		p.logStatRequest(statRequest, "ожидание готовности: ждем ", waitTime.String())
		time.Sleep(waitTime)

		p.logStatRequest(statRequest, "ожидание готовности: попытка ", attempt)

		// Сначала ждем, так как после создания запрос будет собираться

		req, err := p.ozon.StatRequests().Retrieve(statRequest.UUID)
		if err != nil {
			p.logStatRequest(statRequest, "ожидание готовности: ", err)
			continue
		}

		if !req.IsReadyToDownload() {
			p.logStatRequest(statRequest, "ожидание готовности: состояние ", req.State)
			continue
		}

		p.logStatRequest(statRequest, "готов к скачиванию")

		return req, nil
	}

	return ozon.StatRequest{}, errors.New("превышено количество попыток")
}

func (p *processor) downloadStat(statRequest ozon.StatRequest) (string, error) {
	filename := "object-stat-" + statRequest.CampaignId() + ".csv"
	for attempt := 1; attempt <= downloadAttempts; attempt++ {
		p.logStatRequest(statRequest, "скачивание статистики: попытка ", attempt)

		data, err := p.ozon.StatRequests().Download(statRequest)
		if err != nil {
			p.logStatRequest(statRequest, "скачивание статистики: ", err)
			time.Sleep(downloadWaitTime)
			continue
		}

		err = p.storage.Downloads().Write(filename, data)
		if err != nil {
			p.logStatRequest(statRequest, "скачивание статистики: ", err)
			time.Sleep(downloadWaitTime)
			continue
		}

		p.logStatRequest(statRequest, "скачан файл: ", filename)
		return filename, nil
	}

	return "", errors.New("превышено количество попыток скачивания")
}

func (p *processor) logCampaign(c ozon.Campaign, msg ...any) {
	p.logString(
		fmt.Sprintf("[%s] %s\n", c.ID, fmt.Sprint(msg...)),
	)
}

func (p *processor) logStatRequest(r ozon.StatRequest, msg ...any) {
	p.logString(
		fmt.Sprintf("[%s] %s\n", r.CampaignId(), fmt.Sprint(msg...)),
	)
}

func (p *processor) logString(msg string) {
	p.out.Write([]byte(msg))
}
