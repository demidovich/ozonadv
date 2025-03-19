package application

import (
	"ozonadv/config"
	"ozonadv/internal/ozon"
	"ozonadv/internal/stat"
	"ozonadv/internal/storage"
)

type Application struct {
	configFile    string
	config        *config.Config
	ozonCLient    *ozon.Client
	storage       *storage.Storage
	statUsecases  *stat.Usecases
	shutdownFuncs []func()
}

func New() *Application {
	return &Application{}
}

func (a *Application) Config() *config.Config {
	if a.config == nil {
		instance := config.NewOrFail("config.yml")
		a.config = &instance
	}

	return a.config
}

func (a *Application) OzonClient() *ozon.Client {
	if a.ozonCLient == nil {
		a.ozonCLient = ozon.NewClient(a.Config().Ozon)
	}

	return a.ozonCLient
}

func (a *Application) Storage() *storage.Storage {
	if a.storage == nil {
		a.storage = storage.New()
		a.shutdownFuncs = append(a.shutdownFuncs, a.storage.SaveState)
	}

	return a.storage
}

func (a *Application) StatUsecases() *stat.Usecases {
	if a.statUsecases == nil {
		a.statUsecases = stat.New(
			a.Storage(),
			a.OzonClient(),
		)
	}

	return a.statUsecases
}

func (a *Application) Shutdown() {
	for _, f := range a.shutdownFuncs {
		f()
	}
}
