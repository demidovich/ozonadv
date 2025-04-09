package stat1

import (
	"io"
	"ozonadv/internal/ozon"
	"ozonadv/pkg/validation"
	"time"

	"github.com/google/uuid"
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
	Items            []statItem  `json:"items"`
	CreatedAt        string      `json:"createdAt"`
	ApiRequestsCount int         `json:"apiRequestsCount"`
}

func New(options StatOptions) (*Stat, error) {
	if err := options.Validate(); err != nil {
		return nil, err
	}

	instance := &Stat{
		UUID:      uuid.NewString(),
		Options:   options,
		CreatedAt: time.Now().String(),
	}

	return instance, nil
}

func (s *Stat) AddCampaign(campaign ozon.Campaign) {
	for _, i := range s.Items {
		if i.Campaign.ID == campaign.ID {
			return
		}
	}

	s.Items = append(s.Items, statItem{Campaign: campaign})
}

func (s *Stat) ItemByRequestUUID(uuid string) (*statItem, bool) {
	for _, i := range s.Items {
		if i.Request.UUID == uuid {
			return &i, true
		}
	}

	return nil, false
}

func (s *Stat) Start(out io.Writer) {
	o := ozon.New(
		ozon.Config{
			ClientId:     s.Options.CabinetClientId,
			ClientSecret: s.Options.CabinetClientSecret,
		},
		false,
	)

	proc := newProcessor(out, s, o, nil)
	proc.Start()
}

// func (s *stat) ExportToFile(file string) error {
// 	return nil
// }

// func (s *stat) ExportToGoogleSheet() error {
// 	return nil
// }
