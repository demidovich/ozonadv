package application

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
	ozonApi       *ozon.Api
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
		instance := config.NewOrFail("config.yml")
		a.config = &instance
	}

	return a.config
}

func (a *Application) OzonApi() *ozon.Api {
	if a.ozonApi == nil {
		fmt.Println("Инициализация клиента API Озон")
		a.ozonApi = ozon.NewApi(a.Config().Ozon)
	}

	return a.ozonApi
}

func (a *Application) Storage() *storage.Storage {
	if a.storage == nil {
		fmt.Println("Инициализация локального хранилища")
		a.storage = storage.New()
		a.shutdownFuncs = append(a.shutdownFuncs, a.storage.SaveState)
		fmt.Println("Директория локального хранилища", a.storage.RootDir())
	}

	return a.storage
}

func (a *Application) StatUsecases() *stat.Usecases {
	if a.statUsecases == nil {
		a.statUsecases = stat.New(
			a.Storage(),
			a.OzonApi(),
		)
	}

	return a.statUsecases
}

func (a *Application) FindUsecases() *find.Usecases {
	if a.findUsecases == nil {
		a.findUsecases = find.New(
			a.OzonApi(),
		)
	}

	return a.findUsecases
}

func (a *Application) Shutdown() {
	for _, f := range a.shutdownFuncs {
		f()
	}
}
