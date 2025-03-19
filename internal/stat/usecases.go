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
	pull := pullUsecase{ozonClient: ozonClient, storage: storage}

	return &Usecases{
		storage:       storage,
		ozonClient:    ozonClient,
		createUsecase: createUsecase{pullUsecase: pull, ozonClient: ozonClient, storage: storage},
		infoUsecase:   infoUsecase{ozonClient: ozonClient, storage: storage},
		pullUsecase:   pull,
	}
}

func (u *Usecases) HasStatistics() bool {
	return u.storage.StatisticsSize() > 0
}

func (u *Usecases) RemoveAllStatistics() {
	u.storage.RemoveAllStatistics()
}

func (u *Usecases) Create(options CreateOptions) error {
	return u.createUsecase.Handle(options)
}

func (u *Usecases) Info() error {
	return u.infoUsecase.Handle()
}

func (u *Usecases) Pull(optons PullOptions) error {
	return u.pullUsecase.Handle(optons)
}
