package ozon

type CampaignObject struct {
	ID string `json:"id"`
}

func (c *Client) CampaignObjects(campaignId string) ([]CampaignObject, error) {
	result := struct {
		List []CampaignObject `json:"list"`
	}{}

	err := c.Get("/client/campaign/"+campaignId+"/objects", &result)

	return result.List, err
}
