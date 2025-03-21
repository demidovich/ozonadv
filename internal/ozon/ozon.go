package ozon

type Config struct {
	ClientId     string
	ClientSecret string
}

type Ozon struct {
	Campaigns    campaigns
	StatRequests statRequests
}

func New(config Config) *Ozon {
	api := newApi(config)
	return &Ozon{
		Campaigns:    campaigns{api: api},
		StatRequests: statRequests{api: api},
	}
}
