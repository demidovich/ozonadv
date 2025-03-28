package object_stat

import (
	"ozonadv/internal/ozon"
	"ozonadv/internal/storage"
)

type Usecases struct {
	storage    *storage.Storage
	stat       statUsecase
	statInfo   statInfoUsecase
	statExport statExportUsecase
	statReset  statResetUsecase
}

func New(storage *storage.Storage, ozon *ozon.Ozon) *Usecases {
	return &Usecases{
		storage:    storage,
		stat:       statUsecase{ozon: ozon, storage: storage},
		statInfo:   statInfoUsecase{storage: storage},
		statExport: statExportUsecase{storage: storage},
		statReset:  statResetUsecase{storage: storage},
	}
}

func (u *Usecases) HasProcessing() bool {
	return u.storage.StatCampaigns().Size() > 0
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

func (u *Usecases) StatExport(options StatExportOptions) error {
	return u.statExport.Handle(options)
}

func (u *Usecases) StatReset() error {
	return u.statReset.Handle()
}
