package ozon

import (
	"ozonadv/pkg/validation"
)

type statRequests struct {
	api *api
}

func (s *statRequests) All() ([]StatRequest, error) {
	type item struct {
		Meta StatRequest `json:"meta"`
	}

	response := struct {
		Items []item `json:"items"`
		Total string `json:"total"`
	}{}

	err := s.api.httpGet("/client/statistics/externallist", &response)

	result := []StatRequest{}
	for _, item := range response.Items {
		result = append(result, item.Meta)
	}

	return result, err
}

type CreateStatRequestOptions struct {
	CampaignId string `validate:"required,numeric"`
	DateFrom   string `validate:"required,datetime=2006-01-02"`
	DateTo     string `validate:"required,datetime=2006-01-02"`
	GroupBy    string `validate:"required,oneof=NO_GROUP_BY DATE START_OF_WEEK START_OF_MONTH"`
}

func (s *CreateStatRequestOptions) validate() error {
	return validation.ValidateStruct(s)
}

func (s *statRequests) Create(campaign Campaign, options CreateStatRequestOptions) (StatRequest, error) {
	if err := options.validate(); err != nil {
		return StatRequest{}, err
	}

	resource := "/client/statistics/json"
	if campaign.AdvObjectType == "VIDEO_BANNER" {
		resource = "/client/statistics/video/json"
	}

	payload := map[string]any{
		"campaigns": []string{campaign.ID},
		"dateFrom":  options.DateFrom,
		"dateTo":    options.DateTo,
		"groupBy":   options.GroupBy,
	}

	response := struct {
		UUID   string `json:"UUID"`
		Vendor bool   `json:"vendor"`
	}{}

	err := s.api.httpPost(resource, payload, &response)
	if err != nil {
		return StatRequest{}, err
	}

	statRequest := StatRequest{}
	statRequest.UUID = response.UUID
	statRequest.Request.CampaignId = campaign.ID

	return statRequest, nil
}

func (s *statRequests) Retrieve(uuid string) (StatRequest, error) {
	resource := "/client/statistics/" + uuid
	response := StatRequest{}

	err := s.api.httpGet(resource, &response)
	if err != nil {
		return StatRequest{}, err
	}

	return response, nil
}

func (s *statRequests) Download(statRequest StatRequest) ([]byte, error) {
	data, err := s.api.httpGetRaw(statRequest.Link)
	if err != nil {
		return nil, err
	}

	return data, nil
}
