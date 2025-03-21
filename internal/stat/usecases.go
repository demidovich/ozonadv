package stat

import (
	"ozonadv/internal/ozon"
	"ozonadv/internal/storage"
)

type Usecases struct {
	storage          *storage.Storage
	ozonApi          *ozon.Api
	statUsecase      statUsecase
	statInfoUsecase  statInfoUsecase
	statResetUsecase statResetUsecase
}

func New(storage *storage.Storage, ozonApi *ozon.Api) *Usecases {
	return &Usecases{
		storage:          storage,
		ozonApi:          ozonApi,
		statUsecase:      statUsecase{ozonApi: ozonApi, storage: storage},
		statInfoUsecase:  statInfoUsecase{storage: storage},
		statResetUsecase: statResetUsecase{storage: storage},
	}
}

func (u *Usecases) HasIncompleteProcessing() bool {
	return u.storage.Campaigns.Size() > 0
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
