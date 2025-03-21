package find

import (
	"ozonadv/internal/ozon"
)

type Usecases struct {
	findCampaigns findCampaignsUsecase
}

func New(ozon *ozon.Ozon) *Usecases {
	return &Usecases{
		findCampaigns: findCampaignsUsecase{ozon: ozon},
	}
}

func (u *Usecases) Campaigns() error {
	return u.findCampaigns.Handle()
}
