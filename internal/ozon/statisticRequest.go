package ozon

import (
	"ozonadv/pkg/validation"
)

// StatisticRequest State
//
// NOT_STARTED — запрос ожидает выполнения;
// IN_PROGRESS — запрос выполняется в данный момент;
// ERROR       — выполнение запроса завершилось ошибкой;
// OK 		   — запрос успешно выполнен.

// Statistic Group
//
// NO_GROUP_BY
// DATE
// START_OF_WEEK
// START_OF_MONTH

type StatisticRequest struct {
	UUID      string `json:"uuid"`
	State     string `json:"state"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Error     string `json:"error"`
	Link      string `json:"link"`
	Request   struct {
		CampaignId string   `json:"campaignId"`
		Campaigns  []string `json:"campaigns"`
		DateFrom   string   `json:"dateFrom"`
		DateTo     string   `json:"dateTo"`
		From       string   `json:"from"`
		To         string   `json:"to"`
		GroupBy    string   `json:"groupBy"`
	}
}

func (s *StatisticRequest) CampaignId() (value string) {
	switch true {
	case len(s.Request.Campaigns) > 0:
		value = s.Request.Campaigns[0]
	case s.Request.CampaignId != "":
		value = s.Request.CampaignId
	}

	return
}

func (s *StatisticRequest) DateFrom() (value string) {
	switch true {
	case s.Request.DateFrom != "":
		value = s.Request.DateFrom
	case s.Request.From != "":
		value = s.Request.From
	}

	return
}

func (s *StatisticRequest) DateTo() (value string) {
	switch true {
	case s.Request.DateTo != "":
		value = s.Request.DateTo
	case s.Request.To != "":
		value = s.Request.To
	}

	return
}

func (s *StatisticRequest) GroupBy() string {
	return s.Request.GroupBy
}

func (s *StatisticRequest) IsReadyToDownload() bool {
	return s.State == "OK"
}

func (a *Api) StatisticRequests() ([]StatisticRequest, error) {
	type item struct {
		Meta StatisticRequest `json:"meta"`
	}

	response := struct {
		Items []item `json:"items"`
		Total string `json:"total"`
	}{}

	err := a.get("/client/statistics/externallist", &response)

	result := []StatisticRequest{}
	for _, item := range response.Items {
		result = append(result, item.Meta)
	}

	return result, err
}

type StatisticRequestOptions struct {
	CampaignId string `validate:"required,numeric"`
	DateFrom   string `validate:"required,datetime=2006-01-02"`
	DateTo     string `validate:"required,datetime=2006-01-02"`
	GroupBy    string `validate:"required,oneof=NO_GROUP_BY DATE START_OF_WEEK START_OF_MONTH"`
}

func (s *StatisticRequestOptions) validate() error {
	return validation.ValidateStruct(s)
}

func (a *Api) CreateStatisticRequest(campaign Campaign, options StatisticRequestOptions) (*StatisticRequest, error) {
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

	err := a.post(resource, payload, &result)
	if err != nil {
		return nil, err
	}

	statRequest := StatisticRequest{}
	statRequest.UUID = result.UUID

	return &statRequest, nil
}
