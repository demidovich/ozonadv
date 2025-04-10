package stat1

import (
	"encoding/json"
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

type stat struct {
	UUID             string      `json:"uuid"`
	Options          StatOptions `json:"options"`
	Items            []statItem  `json:"items"`
	CreatedAt        string      `json:"createdAt"`
	ApiRequestsCount int         `json:"apiRequestsCount"`
}

func New(options StatOptions) (*stat, error) {
	if err := options.Validate(); err != nil {
		return nil, err
	}

	self := &stat{
		UUID:      uuid.NewString(),
		Options:   options,
		CreatedAt: time.Now().String(),
	}

	return self, nil
}

func FromJson(j string) (*stat, error) {
	self := &stat{}
	err := json.Unmarshal([]byte(j), self)

	return self, err
}

func (s *stat) AddCampaign(campaign ozon.Campaign) {
	for _, i := range s.Items {
		if i.Campaign.ID == campaign.ID {
			return
		}
	}

	s.Items = append(s.Items, statItem{Campaign: campaign})
}

func (s *stat) ItemByRequestUUID(uuid string) (*statItem, bool) {
	for _, i := range s.Items {
		if i.Request.UUID == uuid {
			return &i, true
		}
	}

	return nil, false
}

func (s *stat) Start(out io.Writer) {
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

func (s *stat) ToJson() (string, error) {
	j, err := json.Marshal(s)

	return string(j), err
}
