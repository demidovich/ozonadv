package app

import (
	"io"
	"ozonadv/internal/cabinets"
	"ozonadv/internal/infra/storage"
	"ozonadv/internal/stats"
)

type Application struct {
	out             io.Writer
	storage         *storage.Storage
	cabinetsService *cabinets.Service
	statsService    *stats.Service
	shutdownFuncs   []func()
	debug           Debug
}

func New(out io.Writer) *Application {
	return &Application{
		out:   out,
		debug: newDebug(out),
	}
}

func (a *Application) Debug() Debug {
	return a.debug
}

func (a *Application) Storage() *storage.Storage {
	if a.storage == nil {
		a.debug.Println("[app init] локальное хранилище")
		a.storage = storage.NewDefault()
		// a.RegisterShutdownFunc(a.storage.SaveState)
		a.debug.Println("[app init] директория локального хранилища", a.storage.RootDir())
	}

	return a.storage
}

func (a *Application) CabinetsService() *cabinets.Service {
	if a.cabinetsService == nil {
		a.debug.Println("[app init] сервис рекламных кабинетов")
		a.cabinetsService = cabinets.NewService(a.out, a.Storage().Cabinets(), a.debug)
	}

	return a.cabinetsService
}

func (a *Application) StatsService() *stats.Service {
	if a.statsService == nil {
		a.debug.Println("[app init] сервис статистики")
		a.statsService = stats.NewService(a.Storage().Stats(), a.debug)
	}

	return a.statsService
}

func (a *Application) RegisterShutdownFunc(f func()) {
	a.shutdownFuncs = append(a.shutdownFuncs, f)
}

func (a *Application) Shutdown() {
	a.debug.Println("[app down]")
	for _, f := range a.shutdownFuncs {
		f()
	}
}
