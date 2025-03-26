package stat_processor

import (
	"errors"
	"fmt"
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

type StatProcessor struct {
	storage     *storage.Storage
	ozon        *ozon.Ozon
	procresults *procresults
}

func New(o *ozon.Ozon, s *storage.Storage) *StatProcessor {
	return &StatProcessor{
		ozon:        o,
		storage:     s,
		procresults: newProcresults(),
	}
}

func (p *StatProcessor) Start(c []ozon.Campaign) <-chan bool {
	for _, campaign := range c {
		p.procresults.RegisterCampaign(campaign)
	}

	campaigns := make(chan ozon.Campaign)
	go func() {
		defer close(campaigns)
		for _, campaign := range c {
			campaigns <- campaign
		}
	}()

	createdRequests := p.createStatRequestsStage(campaigns)
	readyRequests := p.readyStatRequestsStage(createdRequests)
	complete := p.downloadStatsStage(readyRequests)

	return complete
}

func (p *StatProcessor) PrintSummaryTable() {
	p.procresults.PrintSummaryTable()
	fmt.Println("")
}

func (p *StatProcessor) createStatRequestsStage(in <-chan ozon.Campaign) <-chan ozon.StatRequest {
	out := make(chan ozon.StatRequest)

	go func() {
		defer close(out)
		for campaign := range in {
			statRequest := p.createStatRequest(campaign)
			out <- statRequest
		}
	}()

	return out
}

func (p *StatProcessor) readyStatRequestsStage(in <-chan ozon.StatRequest) <-chan ozon.StatRequest {
	out := make(chan ozon.StatRequest)

	go func() {
		defer close(out)
		for r := range in {
			statRequest, err := p.readyStatRequest(r)
			if err == nil {
				out <- statRequest
			} else {
				logStatRequest(r, err)
			}
			p.procresults.RegisterStatRequest(statRequest)
		}
	}()

	return out
}

func (p *StatProcessor) downloadStatsStage(in <-chan ozon.StatRequest) <-chan bool {
	out := make(chan bool)

	go func() {
		defer close(out)
		for statRequest := range in {
			filename, err := p.downloadStat(statRequest)
			if err == nil {
				p.storage.Campaigns().Remove(statRequest.Request.CampaignId)
			} else {
				logStatRequest(statRequest, err)
			}
			p.procresults.RegisterStatDownload(statRequest, filename)
		}
		out <- true
	}()

	return out
}

func (p *StatProcessor) createStatRequest(campaign ozon.Campaign) ozon.StatRequest {
	options := ozon.CreateStatRequestOptions{
		CampaignId: campaign.ID,
		DateFrom:   p.storage.StatOptions().DateFrom,
		DateTo:     p.storage.StatOptions().DateTo,
		GroupBy:    p.storage.StatOptions().GroupBy,
	}

	for attempt := 1; attempt <= createAttempts; attempt++ {
		logCampaign(campaign, "создание запроса отчета: попытка ", attempt)

		req, err := p.ozon.StatRequests().Create(campaign, options)
		if err == nil {
			logCampaign(campaign, "создан запрос ", req.UUID)
			return req
		}

		logCampaign(campaign, err)

		waitTime := createWaitTime + createWaitTime*time.Duration(attempt-1)
		logCampaign(campaign, "создание запроса отчета: ждем ", waitTime.String())
		time.Sleep(waitTime)
	}

	fmt.Println("Превышено количество количество попыток создания запроса.")
	fmt.Println("Возможно существует тяжелый несформированый запрос.")
	fmt.Println("Пока Озон не закончит его формирование создать новый не получится.")
	fmt.Println("")
	os.Exit(1)

	return ozon.StatRequest{}
}

func (p *StatProcessor) readyStatRequest(statRequest ozon.StatRequest) (ozon.StatRequest, error) {
	for attempt := 1; attempt <= readyAttempts; attempt++ {
		waitTime := readyWaitTime + readyWaitTime*time.Duration(attempt-1)
		logStatRequest(statRequest, "ожидание готовности: ждем ", waitTime.String())
		time.Sleep(waitTime)

		logStatRequest(statRequest, "ожидание готовности: попытка ", attempt)

		// Сначала ждем, так как после создания запрос будет собираться

		req, err := p.ozon.StatRequests().Retrieve(statRequest.UUID)
		if err != nil {
			logStatRequest(statRequest, "ожидание готовности:", err)
			continue
		}

		if !req.IsReadyToDownload() {
			logStatRequest(statRequest, "ожидание готовности: состояние ", req.State)
			continue
		}

		logStatRequest(statRequest, "готов к скачиванию")

		return req, nil
	}

	return ozon.StatRequest{}, errors.New("превышено количество попыток")
}

func (p *StatProcessor) downloadStat(statRequest ozon.StatRequest) (string, error) {
	filename := statRequest.UUID + ".json"

	for attempt := 1; attempt <= downloadAttempts; attempt++ {
		logStatRequest(statRequest, "скачивание статистики: попытка ", attempt)

		data, err := p.ozon.StatRequests().Download(statRequest)
		if err != nil {
			logStatRequest(statRequest, "скачивание статистики: ", err)
			time.Sleep(downloadWaitTime)
			continue
		}

		err = p.storage.Downloads().Write(filename, data)
		if err != nil {
			logStatRequest(statRequest, "скачивание статистики: ", err)
			time.Sleep(downloadWaitTime)
			continue
		}

		logStatRequest(statRequest, "скачан файл:", filename)
	}

	return "", errors.New("превышено количество попыток скачивания")
}

func logCampaign(c ozon.Campaign, msg ...any) {
	fmt.Printf("[%s] %s\n", c.ID, fmt.Sprint(msg...))
}

func logStatRequest(r ozon.StatRequest, msg ...any) {
	fmt.Printf("[%s] %s\n", r.Request.CampaignId, fmt.Sprint(msg...))
}
