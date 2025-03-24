package stat

import (
	"ozonadv/internal/ozon"
	"ozonadv/internal/storage"
)

type Usecases struct {
	storage   *storage.Storage
	stat      statUsecase
	statInfo  statInfoUsecase
	statReset statResetUsecase
}

func New(storage *storage.Storage, ozon *ozon.Ozon) *Usecases {
	return &Usecases{
		storage:   storage,
		stat:      statUsecase{ozon: ozon, storage: storage},
		statInfo:  statInfoUsecase{storage: storage},
		statReset: statResetUsecase{storage: storage},
	}
}

func (u *Usecases) HasIncompleteProcessing() bool {
	return u.storage.Campaigns().Size() > 0
}

func (u *Usecases) StatNew(options StatOptions) error {
	return u.stat.HandleNew(options)
}

func (u *Usecases) StatContinue() error {
	return u.stat.HandleContinue()
}

func (u *Usecases) StatInfo() error {
	return u.statInfo.Handle()
}

func (u *Usecases) StatReset() error {
	return u.statReset.Handle()
}
