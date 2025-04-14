package ui

import (
	"ozonadv/internal/stats"
)

type statsPage struct {
	statsService *stats.Service
}

func newStats(statsService *stats.Service) statsPage {
	return statsPage{statsService: statsService}
}

func (c statsPage) Run() error {
	return nil
}
