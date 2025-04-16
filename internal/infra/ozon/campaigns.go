package ozon

import (
	"ozonadv/internal/models"
	"strings"
)

type campaigns struct {
	api *api
}

type FindCampaignsFilters struct {
	Ids []string
}

func (c *campaigns) Find(filters FindCampaignsFilters) ([]models.Campaign, error) {
	params := ""
	if len(filters.Ids) > 0 {
		params = "?campaignIds=" + strings.Join(filters.Ids, ",")
	}

	response := struct {
		List  []models.Campaign `json:"list"`
		Total string            `json:"total"`
	}{}

	url := urlApi("/client/campaign") + params
	err := c.api.httpGet(url, &response)

	return response.List, err
}

func (c *campaigns) All() ([]models.Campaign, error) {
	return c.Find(FindCampaignsFilters{})
}
