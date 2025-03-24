package stat_processor

import (
	"fmt"
	"math/rand"
	"ozonadv/internal/ozon"
	"ozonadv/internal/storage"
	"strconv"
	"time"
)

type StatProcessor struct {
	storage *storage.Storage
	ozon    *ozon.Ozon
}

func Start(o *ozon.Ozon, s *storage.Storage, campaigns <-chan ozon.Campaign) <-chan string {
	proc := StatProcessor{
		ozon:    o,
		storage: s,
	}

	newRequests := proc.newStatRequestsStage(campaigns)
	readyRequests := proc.readyStatRequestsStage(newRequests)
	statFiles := proc.downloadStatsStage(readyRequests)

	return statFiles
}

func (p *StatProcessor) newStatRequestsStage(in <-chan ozon.Campaign) <-chan ozon.StatRequest {
	out := make(chan ozon.StatRequest)

	go func() {
		defer close(out)

		for c := range in {
			fmt.Printf("Кампания #%s: создание запроса формирования отчета\n", c.ID)

			time.Sleep(500 * time.Millisecond)

			r := ozon.StatRequest{
				UUID:  strconv.Itoa(rand.Intn(1000000)),
				State: "NOT_STARTED",
			}
			r.Request.CampaignId = c.ID

			fmt.Printf("Кампания #%s, запрос #%s: создан\n", r.Request.CampaignId, r.UUID)
			out <- r
			// statRequest = dirtyRequest()
			// out <- statRequest
		}
	}()

	return out
}

func (p *StatProcessor) readyStatRequestsStage(in <-chan ozon.StatRequest) <-chan ozon.StatRequest {
	out := make(chan ozon.StatRequest)

	go func() {
		defer close(out)

		for r := range in {
			fmt.Printf("Кампания #%s, запрос #%s: ожидание готовности\n", r.Request.CampaignId, r.UUID)

			time.Sleep(500 * time.Millisecond)

			r.State = "OK"
			r.Link = "https://download.ru/report"

			out <- r
			// statRequest = readyRequest()
			// out <- statRequest
		}
	}()

	return out
}

func (p *StatProcessor) downloadStatsStage(in <-chan ozon.StatRequest) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)

		for r := range in {
			fmt.Printf("Кампания #%s, запрос #%s: скачивание статистики\n", r.Request.CampaignId, r.UUID)

			time.Sleep(500 * time.Millisecond)

			p.storage.Campaigns.Remove(r.Request.CampaignId)
			out <- r.Link
			// filename = downloadStat()
			// out <- statRequest
		}
	}()

	return out
}
