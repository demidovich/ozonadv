package app

import (
	"fmt"
	"io"
	"ozonadv/internal/cabinets"
	"ozonadv/internal/stats"
	"ozonadv/internal/storage"
)

type Application struct {
	out             io.Writer
	storage         *storage.Storage
	cabinetsService *cabinets.Service
	statsService    *stats.Service
	shutdownFuncs   []func()
}

func New(out io.Writer) *Application {
	return &Application{
		out: out,
	}
}

func (a *Application) Storage() *storage.Storage {
	if a.storage == nil {
		a.println("[app init] локальное хранилище")
		a.storage = storage.NewTemp()
		// a.RegisterShutdownFunc(a.storage.SaveState)
		a.println("[app init] директория локального хранилища", a.storage.RootDir())
	}

	return a.storage
}

func (a *Application) CabinetsService() *cabinets.Service {
	if a.cabinetsService == nil {
		a.println("[app init] сервис рекламных кабинетов")
		a.cabinetsService = cabinets.NewService(a.out, a.Storage().Cabinets())
	}

	return a.cabinetsService
}

func (a *Application) StatsService() *stats.Service {
	if a.statsService == nil {
		a.println("[app init] сервис статистики")
		a.statsService = stats.NewService(a.out, a.Storage().Stats())
	}

	return a.statsService
}

func (a *Application) RegisterShutdownFunc(f func()) {
	a.shutdownFuncs = append(a.shutdownFuncs, f)
}

func (a *Application) Shutdown() {
	a.println("[app down]")
	for _, f := range a.shutdownFuncs {
		f()
	}
}

func (a *Application) println(m ...any) {
	fmt.Fprintln(a.out, m...)
}

func (a *Application) printf(format string, m ...any) {
	fmt.Fprintf(a.out, format, m...)
}
