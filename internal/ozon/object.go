package ozon

type CampaignObject struct {
	ID string `json:"id"`
}

func (a *Api) CampaignObjects(campaignId string) ([]CampaignObject, error) {
	result := struct {
		List []CampaignObject `json:"list"`
	}{}

	err := a.get("/client/campaign/"+campaignId+"/objects", &result)

	return result.List, err
}
