package ozon

// Statistic State
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

type Statistic struct {
	UUID      string `json:"uuid"`
	State     string `json:"state"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Error     string `json:"error"`
	Link      string `json:"link"`
}

func (s *Statistic) IsReadyToDownload() bool {
	return s.State == "OK"
}

func (c *Client) Statistics() ([]Statistic, error) {
	type item struct {
		Meta Statistic `json:"meta"`
	}

	response := struct {
		Items []item `json:"items"`
		Total string `json:"total"`
	}{}

	err := c.Get("/client/statistics/externallist", &response)

	result := []Statistic{}
	for _, item := range response.Items {
		result = append(result, item.Meta)
	}

	return result, err
}
