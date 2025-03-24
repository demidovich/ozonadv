package ozon

type Config struct {
	ClientId     string
	ClientSecret string
}

type Ozon struct {
	api          *api
	campaigns    *campaigns
	statRequests *statRequests
}

func New(config Config, verbose bool) *Ozon {
	api := newApi(config, verbose)
	return &Ozon{
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

func (o *Ozon) ApiRequestsCount() int {
	return o.api.RequestsCount()
}
