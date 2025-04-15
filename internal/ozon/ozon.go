package ozon

type Config struct {
	ClientId     string
	ClientSecret string
}

type Ozon struct {
	debug        Debug
	api          *api
	campaigns    *campaigns
	statRequests *statRequests
}

func New(config Config, debug Debug) *Ozon {
	api := newApi(config, debug)
	return &Ozon{
		debug:        debug,
		api:          api,
		campaigns:    &campaigns{api: api},
		statRequests: &statRequests{api: api},
	}
}

func (o *Ozon) Campaigns() *campaigns {
	return o.campaigns
}

func (o *Ozon) StatRequests() *statRequests {
	return o.statRequests
}

func (o *Ozon) ApiUsageInfo() {
	o.debug.Println("[shutdown] выполнено запросов ozon api:", o.api.RequestsCount())
}
