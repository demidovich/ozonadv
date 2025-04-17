package export

type Stat map[string]struct {
	Title  string `json:"title"`
	Report struct {
		Rows   []StatRow `json:"rows"`
		Totals struct {
			MoneySpent string `json:"moneySpent"`
			BonusSpent string `json:"bonusSpent"`
		} `json:"totals"`
	} `json:"report"`
}

type StatRow struct {
	CampaignId         string `json:"campaignId"` // Кампания
	CampaignType       string `json:"campaignType"`
	CampaignTitle      string `json:"campaignTitle"`
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

func statRowCsvHeaders() []string {
	return []string{
		"СampaignId",
		"СampaignType",
		"СampaignTitle",
		"Date",
		"Views",
		"Clicks",
		"Ctr",
		"Reach",
		"MoneySpent",
		"Banner",
		"PageType",
		"SearchQuery",
		"AttributesPlatform",
		"AvgBid",
		"IabViews",
		"IabViewsRatio",
		"ViewsQ25",
		"ViewsQ50",
		"ViewsQ75",
		"ViewsQ100",
		"ViewsQ25Ratio",
		"ViewsQ50Ratio",
		"ViewsQ75Ratio",
		"ViewsQ100Ratio",
		"ViewsWithSound",
		"Orders",
	}
}

func statRowCsvValues(s StatRow) []string {
	return []string{
		s.CampaignId,
		s.CampaignType,
		s.CampaignTitle,
		s.Date,
		s.Views,
		s.Clicks,
		s.Ctr,
		s.Reach,
		s.MoneySpent,
		s.Banner,
		s.PageType,
		s.SearchQuery,
		s.AttributesPlatform,
		s.AvgBid,
		s.IabViews,
		s.IabViewsRatio,
		s.ViewsQ25,
		s.ViewsQ50,
		s.ViewsQ75,
		s.ViewsQ100,
		s.ViewsQ25Ratio,
		s.ViewsQ50Ratio,
		s.ViewsQ75Ratio,
		s.ViewsQ100Ratio,
		s.ViewsWithSound,
		s.Orders,
	}
}
