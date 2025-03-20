package stat

import (
	"ozonadv/internal/ozon"
	"ozonadv/internal/storage"
)

type Usecases struct {
	storage          *storage.Storage
	ozonClient       *ozon.Client
	statUsecase      statUsecase
	statInfoUsecase  statInfoUsecase
	statResetUsecase statResetUsecase
}

func New(storage *storage.Storage, ozonClient *ozon.Client) *Usecases {
	return &Usecases{
		storage:          storage,
		ozonClient:       ozonClient,
		statUsecase:      statUsecase{ozonClient: ozonClient, storage: storage},
		statInfoUsecase:  statInfoUsecase{ozonClient: ozonClient, storage: storage},
		statResetUsecase: statResetUsecase{ozonClient: ozonClient, storage: storage},
	}
}

func (u *Usecases) HasCampaignRequests() bool {
	return u.storage.CampaignRequestsSize() > 0
}

func (u *Usecases) RemoveAllCampaignRequests() {
	u.storage.Reset()
}

func (u *Usecases) StatNew(options StatOptions) error {
	return u.statUsecase.HandleNew(options)
}

func (u *Usecases) StatContinue() error {
	return u.statUsecase.HandleContinue()
}

func (u *Usecases) StatInfo() error {
	return u.statInfoUsecase.Handle()
}

func (u *Usecases) StatReset() error {
	return u.statResetUsecase.Handle()
}
