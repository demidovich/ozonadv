package models

type StatItem struct {
	Campaign Campaign `json:"campaign"`
	Request  struct {
		UUID string `json:"uuid"`
		Link string `json:"link"`
		File string `json:"file"`
	} `json:"request"`
}

func (s *StatItem) State() string {
	var val string

	switch {
	case s.Request.File != "":
		val = "Файл скачан"
	case s.Request.Link != "":
		val = "Готов к скачиванию"
	case s.Request.UUID != "":
		val = "Запрос создан"
	default:
		val = ""
	}

	return val
}
