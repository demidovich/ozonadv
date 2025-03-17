package ozon

// states
// CAMPAIGN_STATE_ARCHIVED
// CAMPAIGN_STATE_FINISHED
// CAMPAIGN_STATE_INACTIVE
// CAMPAIGN_STATE_RUNNING

// advObjectTypes
// GLOBAL_PROMO
// VIDEO_BANNER
// BANNER

type Campaign struct {
	ID                       string `json:"id"`
	Title                    string `json:"title"`
	State                    string `json:"state"`
	AdvObjectType            string `json:"advObjectType"`
	FromDate                 string `json:"fromDate"`
	ToDate                   string `json:"toDate"`
	DailyBudget              string `json:"dailyBudget"`
	Budget                   string `json:"budget"`
	CreatedAt                string `json:"createdAt"`
	UpdatedAt                string `json:"updatedAt"`
	ProductCampaignMode      string `json:"productCampaignMode"`
	ProductAutopilotStrategy string `json:"productAutopilotStrategy"`
}

func (c *Client) Campaigns() ([]Campaign, error) {
	result := struct {
		List  []Campaign `json:"list"`
		Total string     `json:"total"`
	}{}

	err := c.get("/client/campaign", &result)

	return result.List, err
}
