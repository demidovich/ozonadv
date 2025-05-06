package models

import (
	"encoding/json"
	"unicode/utf8"

	"github.com/demidovich/ozonadv/pkg/validation"
)

type StatOptions struct {
	Name                string `validate:"omitempty" json:"name"`
	CabinetUUID         string `validate:"required" json:"cabinetUUID"`
	CabinetName         string `validate:"required" json:"cabinetName"`
	CabinetClientID     string `validate:"required" json:"cabinetClientID"`
	CabinetClientSecret string `validate:"required" json:"cabinetClientSecret"`
	Type                string `validate:"required,oneof=TOTAL OBJECT" json:"type"`
	DateFrom            string `validate:"required,datetime=2006-01-02" json:"dateFrom"`
	DateTo              string `validate:"required,datetime=2006-01-02" json:"dateTo"`
	GroupBy             string `validate:"required,oneof=NO_GROUP_BY DATE START_OF_WEEK START_OF_MONTH" json:"groupBy"`
}

func (s StatOptions) Validate() error {
	return validation.ValidateStruct(s)
}

type Stat struct {
	UUID             string      `json:"uuid"`
	Options          StatOptions `json:"options"`
	Items            []StatItem  `json:"items"`
	RunnedAt         string      `json:"runnedAt"`
	CreatedAt        string      `json:"createdAt"`
	APIRequestsCount int         `json:"apiRequestsCount"`
}

func (s *Stat) StateHuman() string {
	if s.RunnedAt == "" {
		return "Не запускался"
	}

	completed := true
	for _, i := range s.Items {
		if i.Request.Link == "" {
			completed = false
			break
		}
	}

	if completed {
		return "Готов к экспорту"
	} else {
		return "В процессе"
	}
}

func (s *Stat) AddCampaign(campaign Campaign) {
	for _, i := range s.Items {
		if i.Campaign.ID == campaign.ID {
			return
		}
	}

	s.Items = append(s.Items, StatItem{Campaign: campaign})
}

func (s *Stat) Campaigns() []*Campaign {
	result := make([]*Campaign, 0, len(s.Items))
	for i := range s.Items {
		result = append(result, &s.Items[i].Campaign)
	}

	return result
}

func (s *Stat) CampaignsNotStarted() []*Campaign {
	filter := func(i StatItem) bool {
		return i.NotStarted()
	}

	return s.campaignsFiltered(filter)
}

func (s *Stat) CampaignsInProccess() []*Campaign {
	filter := func(i StatItem) bool {
		return i.InProccess()
	}

	return s.campaignsFiltered(filter)
}

func (s *Stat) CampaignsCompleted() []*Campaign {
	filter := func(i StatItem) bool {
		return i.Completed()
	}

	return s.campaignsFiltered(filter)
}

func (s *Stat) campaignsFiltered(filter func(i StatItem) bool) []*Campaign {
	result := make([]*Campaign, 0)
	for i, item := range s.Items {
		if filter(item) {
			result = append(result, &s.Items[i].Campaign)
		}
	}

	return result
}

func (s *Stat) ItemByRequestUUID(uuid string) (*StatItem, bool) {
	for _, i := range s.Items {
		if i.Request.UUID == uuid {
			return &i, true
		}
	}

	return nil, false
}

func (s *Stat) NameTruncated(length int) string {
	suffix := "..."
	suffixLength := utf8.RuneCountInString(suffix)

	if utf8.RuneCountInString(s.Options.Name)+suffixLength <= length {
		return s.Options.Name
	}

	r := []rune(s.Options.Name)
	return string(r[:length-suffixLength]) + suffix
}

func (s *Stat) ToJSON() (string, error) {
	j, err := json.Marshal(s)

	return string(j), err
}

func (s *Stat) IsTypeObject() bool {
	return s.Options.Type == "OBJECT"
}

func (s *Stat) IsTypeTotal() bool {
	return s.Options.Type == "TOTAL"
}
