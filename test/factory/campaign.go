package factory

import (
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/demidovich/ozonadv/internal/models"
)

type campaignFactory struct {
	state    string
	fromDate string
	toDate   string
}

func Campaign() *campaignFactory {
	return &campaignFactory{}
}

func (f *campaignFactory) New() models.Campaign {
	if f.state == "" {
		f.WithStateFinished()
	}

	if f.fromDate == "" {
		f.WithFromDate(time.Now())
	}

	if f.toDate == "" {
		f.WithToDate(time.Now().Add(10 * 24 * time.Hour))
	}

	c := models.Campaign{
		ID:                       strconv.Itoa(gofakeit.Int()),
		Title:                    gofakeit.Company(),
		State:                    f.state,
		AdvObjectType:            "BANNER",
		FromDate:                 f.fromDate,
		ToDate:                   f.toDate,
		DailyBudget:              "504000000",
		Budget:                   "50000000",
		CreatedAt:                time.Now().Add(-15 * 24 * time.Hour).String(),
		UpdatedAt:                time.Now().Add(-10 * 24 * time.Hour).String(),
		ProductCampaignMode:      "PRODUCT_CAMPAIGN_MODE_AUTO",
		ProductAutopilotStrategy: "NO_AUTO_STRATEGY",
	}

	f.reset()

	return c
}

func (f *campaignFactory) WithFromDate(d time.Time) *campaignFactory {
	f.fromDate = d.Format("2006-01-02")
	return f
}

func (f *campaignFactory) WithToDate(d time.Time) *campaignFactory {
	f.toDate = d.Format("2006-01-02")
	return f
}

func (f *campaignFactory) WithStateRunning() *campaignFactory {
	f.state = "CAMPAIGN_STATE_RUNNING"
	return f
}

func (f *campaignFactory) WithStatePlanned() *campaignFactory {
	f.state = "CAMPAIGN_STATE_PLANNED"
	return f
}

func (f *campaignFactory) WithStateStopped() *campaignFactory {
	f.state = "CAMPAIGN_STATE_STOPPED"
	return f
}

func (f *campaignFactory) WithStateInactive() *campaignFactory {
	f.state = "CAMPAIGN_STATE_INACTIVE"
	return f
}

func (f *campaignFactory) WithStateArchived() *campaignFactory {
	f.state = "CAMPAIGN_STATE_ARCHIVED"
	return f
}

func (f *campaignFactory) WithStateModerationDraft() *campaignFactory {
	f.state = "CAMPAIGN_STATE_MODERATION_DRAFT"
	return f
}

func (f *campaignFactory) WithStateModerationInProgress() *campaignFactory {
	f.state = "CAMPAIGN_STATE_MODERATION_IN_PROGRESS"
	return f
}

func (f *campaignFactory) WithStateModerationFailed() *campaignFactory {
	f.state = "CAMPAIGN_STATE_MODERATION_FAILED"
	return f
}

func (f *campaignFactory) WithStateFinished() *campaignFactory {
	f.state = "CAMPAIGN_STATE_FINISHED"
	return f
}

func (f *campaignFactory) reset() {
	f.state = ""
	f.fromDate = ""
	f.toDate = ""
}
