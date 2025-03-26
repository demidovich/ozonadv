package app

import (
	"fmt"
	"ozonadv/config"
	"ozonadv/internal/find"
	"ozonadv/internal/ozon"
	"ozonadv/internal/stat"
	"ozonadv/internal/storage"
)

type Application struct {
	configFile    string
	config        *config.Config
	ozon          *ozon.Ozon
	storage       *storage.Storage
	findUsecases  *find.Usecases
	statUsecases  *stat.Usecases
	shutdownFuncs []func()
}

func New() *Application {
	return &Application{}
}

func (a *Application) Config() *config.Config {
	if a.config == nil {
		fmt.Println("[app init] конфигурация")
		instance := config.NewOrFail("config.yml")
		a.config = &instance
	}

	return a.config
}

func (a *Application) Ozon() *ozon.Ozon {
	if a.ozon == nil {
		fmt.Println("[app init] клиент ozon")
		a.ozon = ozon.New(a.Config().Ozon, a.Config().Verbose)
		a.RegisterShutdownFunc(a.ozon.ApiUsageInfo)
	}

	return a.ozon
}

func (a *Application) Storage() *storage.Storage {
	if a.storage == nil {
		fmt.Println("[app init] локальное хранилище")
		a.storage = storage.New()
		a.RegisterShutdownFunc(a.storage.SaveState)
		fmt.Println("[app init] директория локального хранилища", a.storage.RootDir())
	}

	return a.storage
}

func (a *Application) StatUsecases() *stat.Usecases {
	if a.statUsecases == nil {
		a.statUsecases = stat.New(
			a.Storage(),
			a.Ozon(),
		)
	}

	return a.statUsecases
}

func (a *Application) FindUsecases() *find.Usecases {
	if a.findUsecases == nil {
		a.findUsecases = find.New(
			a.Ozon(),
		)
	}

	return a.findUsecases
}

func (a *Application) RegisterShutdownFunc(f func()) {
	a.shutdownFuncs = append(a.shutdownFuncs, f)
}

func (a *Application) Shutdown() {
	for _, f := range a.shutdownFuncs {
		f()
	}
}
