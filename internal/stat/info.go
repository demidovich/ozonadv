package stat

import (
	"ozonadv/internal/ozon"
	"ozonadv/internal/storage"
)

type infoUsecase struct {
	storage    *storage.Storage
	ozonClient *ozon.Client
}

func (i *infoUsecase) Handle() error {
	return nil
}
