package ozon

// Statistic states
// NOT_STARTED — запрос ожидает выполнения;
// IN_PROGRESS — запрос выполняется в данный момент;
// ERROR       — выполнение запроса завершилось ошибкой;
// OK 		   — запрос успешно выполнен.

type Statistic struct {
	UUID      string `json:"uuid"`
	State     string `json:"state"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Error     string `json:"error"`
	Link      string `json:"link"`
}

func (s *Statistic) IsOk() bool {
	return s.State == "OK"
}

func (c *Client) Statistics() []Statistic {
	return []Statistic{}
}
