package models

// StatRequest State
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

type StatRequest struct {
	UUID      string `json:"uuid"`
	State     string `json:"state"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Error     string `json:"error"`
	Link      string `json:"link"`
	Request   struct {
		CampaignID string   `json:"campaignId"`
		Campaigns  []string `json:"campaigns"`
		DateFrom   string   `json:"dateFrom"`
		DateTo     string   `json:"dateTo"`
		From       string   `json:"from"`
		To         string   `json:"to"`
		GroupBy    string   `json:"groupBy"`
	} `json:"request"`
}

func (s *StatRequest) CampaignID() (value string) {
	switch {
	case len(s.Request.Campaigns) > 0:
		value = s.Request.Campaigns[0]
	case s.Request.CampaignID != "":
		value = s.Request.CampaignID
	}

	return
}

func (s *StatRequest) DateFrom() (value string) {
	switch {
	case s.Request.DateFrom != "":
		value = s.Request.DateFrom
	case s.Request.From != "":
		value = s.Request.From
	}

	return
}

func (s *StatRequest) DateTo() (value string) {
	switch {
	case s.Request.DateTo != "":
		value = s.Request.DateTo
	case s.Request.To != "":
		value = s.Request.To
	}

	return
}

func (s *StatRequest) GroupBy() string {
	return s.Request.GroupBy
}

func (s *StatRequest) IsReadyToDownload() bool {
	return s.State == "OK"
}
