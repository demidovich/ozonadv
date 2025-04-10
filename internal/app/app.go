package app

import (
	"fmt"
	"io"
	"ozonadv/internal/stats"
	"ozonadv/internal/storage"
)

type Application struct {
	out           io.Writer
	storage       *storage.Storage
	statsService  *stats.Service
	shutdownFuncs []func()
}

func New(out io.Writer) *Application {
	return &Application{
		out: out,
	}
}

func (a *Application) Storage() *storage.Storage {
	if a.storage == nil {
		fmt.Println("[app init] локальное хранилище")
		a.storage = storage.NewTemp()
		// a.RegisterShutdownFunc(a.storage.SaveState)
		fmt.Println("[app init] директория локального хранилища", a.storage.RootDir())
	}

	return a.storage
}

func (a *Application) StatsService() *stats.Service {
	if a.statsService == nil {
		fmt.Println("[app init] сервис статистики")
		a.statsService = stats.NewService(a.out, a.Storage().Stats())
	}

	return a.statsService
}

// func (a *Application) CampaignsUsecases() *campaigns.Usecases {
// 	if a.campaignsUsecases == nil {
// 		a.campaignsUsecases = campaigns.New(
// 			a.Storage(),
// 			a.Ozon(),
// 		)
// 	}

// 	return a.campaignsUsecases
// }

// func (a *Application) StatUsecases() *stat.Usecases {
// 	if a.statUsecases == nil {
// 		a.statUsecases = stat.New(
// 			a.Storage(),
// 			a.Ozon(),
// 		)
// 	}

// 	return a.statUsecases
// }

// func (a *Application) ObjectStatUsecases() *object_stat.Usecases {
// 	if a.objectStatUsecases == nil {
// 		a.objectStatUsecases = object_stat.New(
// 			a.Storage(),
// 			a.Ozon(),
// 		)
// 	}

// 	return a.objectStatUsecases
// }

func (a *Application) RegisterShutdownFunc(f func()) {
	a.shutdownFuncs = append(a.shutdownFuncs, f)
}

func (a *Application) Shutdown() {
	for _, f := range a.shutdownFuncs {
		f()
	}
}
