package stat

import (
	"ozonadv/internal/ozon"
	"ozonadv/internal/storage"
)

type pullUsecase struct {
	storage    *storage.Storage
	ozonClient *ozon.Client
}

func (p *pullUsecase) Handle() error {
	return nil
}
