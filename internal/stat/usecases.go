package stat

import (
	"ozonadv/internal/ozon"
	"ozonadv/internal/storage"
)

type Usecases struct {
	storage       *storage.Storage
	ozonClient    *ozon.Client
	createUsecase createUsecase
	infoUsecase   infoUsecase
	pullUsecase   pullUsecase
}

func New(storage *storage.Storage, ozonClient *ozon.Client) *Usecases {
	return &Usecases{
		storage:       storage,
		ozonClient:    ozonClient,
		createUsecase: createUsecase{ozonClient: ozonClient, storage: storage},
		infoUsecase:   infoUsecase{ozonClient: ozonClient, storage: storage},
		pullUsecase:   pullUsecase{ozonClient: ozonClient, storage: storage},
	}
}

func (u *Usecases) HasIncompletedStatistics() bool {
	return u.storage.StatisticsSize() > 0
}

func (u *Usecases) Create(options CreateOptions) error {
	return u.createUsecase.Handle(options)
}

func (u *Usecases) Info() error {
	return u.infoUsecase.Handle()
}

func (u *Usecases) Pull() error {
	return u.pullUsecase.Handle()
}
