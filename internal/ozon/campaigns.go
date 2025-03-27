package ozon

import (
	"strings"
)

type campaigns struct {
	api *api
}

type FindCampaignsFilters struct {
	Ids []string
}

func (c *campaigns) Find(filters FindCampaignsFilters) ([]Campaign, error) {
	params := ""
	if len(filters.Ids) > 0 {
		params = "?campaignIds=" + strings.Join(filters.Ids, ",")
	}

	response := struct {
		List  []Campaign `json:"list"`
		Total string     `json:"total"`
	}{}

	err := c.api.httpGet("/client/campaign"+params, &response)

	return response.List, err
}

func (c *campaigns) All() ([]Campaign, error) {
	return c.Find(FindCampaignsFilters{})
}
