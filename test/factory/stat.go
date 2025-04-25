package factory

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/demidovich/ozonadv/internal/models"
)

type statFactory struct {
	campaignFactory *campaignFactory
	optionsDateFrom string
	optionsDateTo   string
	optionsType     string
	campaigns       []models.Campaign
	campaignsCount  int
}

func Stat() *statFactory {
	return &statFactory{
		campaignFactory: Campaign(),
	}
}

func (f *statFactory) New() *models.Stat {
	if f.optionsDateFrom == "" {
		f.WithDateFrom(time.Now().Add(-10 * 24 * time.Hour))
	}

	if f.optionsDateTo == "" {
		f.WithDateTo(time.Now())
	}

	if f.optionsType == "" {
		f.WithTypeObject()
	}

	options := models.StatOptions{
		Name:                gofakeit.AppName(),
		CabinetUUID:         gofakeit.UUID(),
		CabinetName:         gofakeit.FarmAnimal(),
		CabinetClientID:     gofakeit.UUID(),
		CabinetClientSecret: gofakeit.UUID(),
		Type:                f.optionsType,
		DateFrom:            f.optionsDateFrom,
		DateTo:              f.optionsDateTo,
		GroupBy:             "OBJECT",
	}

	s := models.Stat{
		UUID:    gofakeit.UUID(),
		Options: options,
	}

	if len(f.campaigns) > 0 {
		for _, c := range f.campaigns {
			s.AddCampaign(c)
		}
	} else if f.campaignsCount > 0 {
		for range f.campaignsCount {
			c := f.campaignFactory.New()
			s.AddCampaign(c)
		}
	}

	f.reset()

	return &s
}

func (f *statFactory) WithDateFrom(d time.Time) *statFactory {
	f.optionsDateFrom = d.Format("2006-01-02")
	return f
}

func (f *statFactory) WithDateTo(d time.Time) *statFactory {
	f.optionsDateTo = d.Format("2006-01-02")
	return f
}

func (f *statFactory) WithTypeTotal() *statFactory {
	f.optionsType = "TOTAL"
	return f
}

func (f *statFactory) WithTypeObject() *statFactory {
	f.optionsType = "OBJECT"
	return f
}

func (f *statFactory) WithCampaign(c models.Campaign) *statFactory {
	f.campaigns = append(f.campaigns, c)
	return f
}

func (f *statFactory) WithCampaignsCount(v int) *statFactory {
	if v < 1 {
		return f
	}

	return f
}

func (f *statFactory) reset() {
	f.optionsDateFrom = ""
	f.optionsDateTo = ""
	f.optionsType = ""
	f.campaigns = []models.Campaign{}
	f.campaignsCount = 0
}
