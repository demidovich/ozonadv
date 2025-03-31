package ozon

import (
	"ozonadv/pkg/validation"
	"strings"
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

	url := urlApi("/client/statistics/externallist")
	err := s.api.httpGet(url, &response)

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
	url := urlApi(resource)

	payload := map[string]any{
		"campaigns": []string{campaign.ID},
		"dateFrom":  options.DateFrom,
		"dateTo":    options.DateTo,
		"groupBy":   options.GroupBy,
	}

	response := StatRequest{}
	err := s.api.httpPost(url, payload, &response)

	return response, err
}

func (s *statRequests) CreateObject(campaign Campaign, options CreateStatRequestOptions) (StatRequest, error) {
	if err := options.validate(); err != nil {
		return StatRequest{}, err
	}

	resource := "/statistics"
	if campaign.AdvObjectType == "VIDEO_BANNER" {
		resource = "/statistics/video"
	}
	url := urlAdvApi(resource)

	payload := map[string]any{
		"campaignId":      []string{campaign.ID},
		"dateFrom":        options.DateFrom,
		"dateTo":          options.DateTo,
		"groupBy":         options.GroupBy,
		"attributionDays": "30",
	}

	response := StatRequest{}
	err := s.api.httpPost(url, payload, &response)

	return response, err
}

func (s *statRequests) Retrieve(uuid string) (StatRequest, error) {
	response := StatRequest{}

	url := urlApi("/client/statistics/" + uuid)
	err := s.api.httpGet(url, &response)

	return response, err
}

func (s *statRequests) Download(statRequest StatRequest) ([]byte, error) {
	url := statRequest.Link
	if !strings.HasPrefix(url, "http") {
		url = urlApi(url)
	}

	data, err := s.api.httpGetRaw(url)
	if err != nil {
		return nil, err
	}

	return data, nil
}
