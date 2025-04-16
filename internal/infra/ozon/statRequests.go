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

// Создание общей статистики кампании
func (s *statRequests) CreateTotal(campaign Campaign, options CreateStatRequestOptions) (StatRequest, error) {
	if err := options.validate(); err != nil {
		return StatRequest{}, err
	}

	var resource string
	if campaign.AdvObjectType == "VIDEO_BANNER" {
		resource = "/client/statistics/video/json"
	} else {
		resource = "/client/statistics/json"
	}

	payload := map[string]any{
		"campaigns": []string{campaign.ID},
		"dateFrom":  options.DateFrom,
		"dateTo":    options.DateTo,
		"groupBy":   options.GroupBy,
	}

	response := StatRequest{}
	err := s.api.httpPost(urlApi(resource), payload, &response)

	return response, err
}

// Создание статистики по объектам рекламной кампании
func (s *statRequests) CreateObject(campaign Campaign, options CreateStatRequestOptions) (StatRequest, error) {
	if err := options.validate(); err != nil {
		return StatRequest{}, err
	}

	var resource string
	var payload map[string]any

	if campaign.AdvObjectType == "VIDEO_BANNER" {
		resource = "/statistics/video"
		payload = map[string]any{
			"campaigns": []string{campaign.ID},
			"dateFrom":  options.DateFrom,
			"dateTo":    options.DateTo,
			"groupBy":   options.GroupBy,
		}
	} else {
		resource = "/statistics"
		payload = map[string]any{
			// "campaignId":      campaign.ID,
			"campaigns":       []string{campaign.ID},
			"dateFrom":        options.DateFrom,
			"dateTo":          options.DateTo,
			"groupBy":         options.GroupBy,
			"attributionDays": "30",
		}
	}

	response := StatRequest{}
	err := s.api.httpPost(urlAdvApi(resource), payload, &response)
	response.Request.CampaignId = campaign.ID

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
		url = apiHost + url
	}

	data, err := s.api.httpGetRaw(url)
	if err != nil {
		return nil, err
	}

	return data, nil
}
