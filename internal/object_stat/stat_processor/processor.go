package stat_processor

import (
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

type statProcessor struct {
	storage *storage.Storage
	ozon    *ozon.Ozon
}

func New(o *ozon.Ozon, s *storage.Storage) *statProcessor {
	return &statProcessor{
		ozon:    o,
		storage: s,
	}
}

func (p *statProcessor) Start(c []ozon.Campaign) {

}
