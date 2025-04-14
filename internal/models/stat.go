package models

import (
	"encoding/json"
	"ozonadv/internal/ozon"
	"ozonadv/pkg/validation"
)

type StatOptions struct {
	Name                string `validate:"omitempty" json:"name"`
	CabinetName         string `validate:"required" json:"cabinetName"`
	CabinetClientId     string `validate:"required" json:"cabinetClientID"`
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
	CreatedAt        string      `json:"createdAt"`
	ApiRequestsCount int         `json:"apiRequestsCount"`
}

func (s *Stat) State() string {
	completed := true
	for _, i := range s.Items {
		if i.Request.Link == "" {
			completed = false
			break
		}
	}

	if completed {
		return "Готов"
	} else {
		return "Не готов"
	}
}

func (s *Stat) AddCampaign(campaign ozon.Campaign) {
	for _, i := range s.Items {
		if i.Campaign.ID == campaign.ID {
			return
		}
	}

	s.Items = append(s.Items, StatItem{Campaign: campaign})
}

func (s *Stat) Campaigns() []ozon.Campaign {
	result := make([]ozon.Campaign, 0, len(s.Items))
	for _, i := range s.Items {
		result = append(result, i.Campaign)
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

func (s *Stat) ToJson() (string, error) {
	j, err := json.Marshal(s)

	return string(j), err
}
