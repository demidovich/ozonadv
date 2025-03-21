package ozon

import (
	"fmt"
	"os"
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

	err := s.api.Get("/client/statistics/externallist", &response)

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

func (s *statRequests) Create(campaign Campaign, options CreateStatRequestOptions) (*StatRequest, error) {
	if err := options.validate(); err != nil {
		return nil, err
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

	result := struct {
		UUID   string `json:"UUID"`
		Vendor bool   `json:"vendor"`
	}{}

	err := s.api.Post(resource, payload, &result)
	if err != nil {
		return nil, err
	}

	statRequest := StatRequest{}
	statRequest.UUID = result.UUID

	return &statRequest, nil
}

func (s *statRequests) Retrieve(uuid string) (*StatRequest, error) {
	url := s.api.Url("/client/statistics/" + uuid)

	result := StatRequest{}

	err := s.api.Get(url, &result)
	fmt.Println(url)
	fmt.Println(err)
	os.Exit(1)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *statRequests) Download(statRequest StatRequest) ([]byte, error) {
	url := apiHost + statRequest.Link

	data, err := s.api.GetRaw(url)
	if err != nil {
		return nil, err
	}

	return data, nil
}
