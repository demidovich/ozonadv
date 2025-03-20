package ozon

import "strings"

// Campaign states
// CAMPAIGN_STATE_RUNNING                — активная кампания;
// CAMPAIGN_STATE_PLANNED                — кампания, сроки проведения которой ещё не наступили;
// CAMPAIGN_STATE_STOPPED                — кампания, приостановленная из-за нехватки бюджета;
// CAMPAIGN_STATE_INACTIVE               — кампания, остановленная владельцем;
// CAMPAIGN_STATE_ARCHIVED               — архивная кампания;
// CAMPAIGN_STATE_MODERATION_DRAFT       — отредактированная кампания до отправки на модерацию;
// CAMPAIGN_STATE_MODERATION_IN_PROGRESS — кампания, отправленная на модерацию;
// CAMPAIGN_STATE_MODERATION_FAILED      — кампания, непрошедшая модерацию;
// CAMPAIGN_STATE_FINISHED               — кампания завершена, дата окончания в прошлом, такую кампанию нельзя изменить,
//                                         можно только клонировать или создать новую.

// Campaign advObjectTypes
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

func (c *Campaign) NeverRun() bool {
	return c.State == "CAMPAIGN_STATE_PLANNED"
}

func (c *Campaign) ShortState() string {
	return strings.TrimPrefix(c.State, "CAMPAIGN_STATE_")
}

func (c *Client) AllCampaigns() ([]Campaign, error) {
	response := struct {
		List  []Campaign `json:"list"`
		Total string     `json:"total"`
	}{}

	err := c.get("/client/campaign", &response)

	return response.List, err
}
