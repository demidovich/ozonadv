package stat

type StatExport struct {
	СampaignId         string `json:"campaignId"` // Кампания
	СampaignType       string `json:"campaignType"`
	СampaignTitle      string `json:"campaignTitle"`
	Date               string `json:"date"` // Общие
	Views              string `json:"views"`
	Clicks             string `json:"clicks"`
	Ctr                string `json:"ctr"`
	Reach              string `json:"reach"`
	MoneySpent         string `json:"moneySpent"`
	Banner             string `json:"banner"`
	PageType           string `json:"pageType"` // Обычная реклама
	SearchQuery        string `json:"search_query"`
	AttributesPlatform string `json:"attributesPlatform"`
	AvgBid             string `json:"avgBid"`
	IabViews           string `json:"iabViews"` // Поля видеорекламы
	IabViewsRatio      string `json:"iabViewsRatio"`
	ViewsQ25           string `json:"viewsQ25"`
	ViewsQ50           string `json:"viewsQ50"`
	ViewsQ75           string `json:"viewsQ75"`
	ViewsQ100          string `json:"viewsQ100"`
	ViewsQ25Ratio      string `json:"viewsQ25Ratio"`
	ViewsQ50Ratio      string `json:"viewsQ50Ratio"`
	ViewsQ75Ratio      string `json:"viewsQ75Ratio"`
	ViewsQ100Ratio     string `json:"viewsQ100Ratio"`
	ViewsWithSound     string `json:"viewsWithSound"`
	Orders             string `json:"orders"`
}
