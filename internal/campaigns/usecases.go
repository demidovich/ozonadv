package campaigns

import (
	"ozonadv/internal/ozon"
	"ozonadv/internal/storage"
)

type Usecases struct {
	selectCampaigns selectCampaignsUsecase
}

func New(storage *storage.Storage, ozon *ozon.Ozon) *Usecases {
	return &Usecases{
		selectCampaigns: selectCampaignsUsecase{storage: storage, ozon: ozon},
	}
}

func (u *Usecases) Select() error {
	return u.selectCampaigns.Handle()
}
