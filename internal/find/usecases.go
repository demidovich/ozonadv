package find

import (
	"ozonadv/internal/ozon"
)

type Usecases struct {
	ozonApi              *ozon.Api
	findCampaignsUsecase findCampaignsUsecase
}

func New(ozonApi *ozon.Api) *Usecases {
	return &Usecases{
		ozonApi:              ozonApi,
		findCampaignsUsecase: findCampaignsUsecase{ozonApi: ozonApi},
	}
}

func (u *Usecases) Campaigns() error {
	return u.findCampaignsUsecase.Handle()
}
