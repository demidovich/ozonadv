package find

import (
	"ozonadv/internal/ozon"
	"ozonadv/internal/storage"
)

type Usecases struct {
	findCampaigns findCampaignsUsecase
}

func New(storage *storage.Storage, ozon *ozon.Ozon) *Usecases {
	return &Usecases{
		findCampaigns: findCampaignsUsecase{storage: storage, ozon: ozon},
	}
}

func (u *Usecases) Campaigns() error {
	return u.findCampaigns.Handle()
}
