package ozon

import (
	"ozonadv/internal/models"
	"strings"
)

type campaigns struct {
	api *api
}

type FindCampaignsFilters struct {
	IDs []string
}

func (c *campaigns) Find(filters FindCampaignsFilters) ([]models.Campaign, error) {
	params := ""
	if len(filters.IDs) > 0 {
		params = "?campaignIDs=" + strings.Join(filters.IDs, ",")
	}

	response := struct {
		List  []models.Campaign `json:"list"`
		Total string            `json:"total"`
	}{}

	url := urlAPI("/client/campaign") + params
	err := c.api.httpGet(url, &response)

	return response.List, err
}

func (c *campaigns) All() ([]models.Campaign, error) {
	return c.Find(FindCampaignsFilters{})
}
