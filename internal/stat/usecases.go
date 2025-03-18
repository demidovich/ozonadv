package stat

import (
	"ozonadv/internal/ozon"
	"ozonadv/internal/storage"
)

type Usecases struct {
	storage    *storage.Storage
	ozonClient *ozon.Client
}

func New(storage *storage.Storage, ozonClient *ozon.Client) Usecases {
	return Usecases{
		storage:    storage,
		ozonClient: ozonClient,
	}
}
