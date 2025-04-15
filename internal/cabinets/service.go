package cabinets

import (
	"io"
	"ozonadv/internal/models"
	"ozonadv/internal/ozon"
	"slices"
	"strings"
	"time"
)

type Service struct {
	storage        storage
	campaignsCache []ozon.Campaign
	debug          Debug
}

func NewService(out io.Writer, storage storage, debug Debug) *Service {
	return &Service{
		storage: storage,
		debug:   debug,
	}
}

func (s *Service) All() []models.Cabinet {
	return s.storage.All()
}

func (s *Service) Find(uuid string) (*models.Cabinet, bool) {
	for _, c := range s.storage.All() {
		if c.UUID == uuid {
			return &c, true
		}
	}
	return nil, false
}

func (s *Service) Add(c models.Cabinet) error {
	if err := c.Validate(); err != nil {
		return err
	}

	if !s.storage.Has(c) {
		c.CreatedAt = time.Now().String()
	}

	s.storage.Add(c)

	return nil
}

func (s *Service) Remove(c models.Cabinet) {
	s.storage.Remove(c)
}

type CampaignFilters struct {
	Title  string
	States []string
}

func (s *Service) Campaigns(cabinet models.Cabinet) ([]ozon.Campaign, error) {
	if len(s.campaignsCache) == 0 {
		var err error
		if s.campaignsCache, err = s.ozon(cabinet).Campaigns().All(); err != nil {
			return s.campaignsCache, err
		}
	}

	return s.campaignsCache, nil
}

func (s *Service) CampaignsFiltered(cabinet models.Cabinet, filters CampaignFilters) ([]ozon.Campaign, error) {
	result, err := s.ozon(cabinet).Campaigns().All()
	if err != nil {
		return result, err
	}

	if filters.Title != "" {
		filtered := []ozon.Campaign{}
		for _, campaign := range result {
			if strings.Contains(strings.ToLower(campaign.Title), strings.ToLower(filters.Title)) {
				filtered = append(filtered, campaign)
			}
		}
		result = filtered
	}

	if len(filters.States) > 0 {
		filtered := []ozon.Campaign{}
		for _, campaign := range result {
			if slices.Contains(filters.States, campaign.State) {
				filtered = append(filtered, campaign)
			}
		}
		result = filtered
	}

	return result, nil
}

func (s *Service) ozon(c models.Cabinet) *ozon.Ozon {
	return ozon.New(
		ozon.Config{
			ClientId:     c.ClientID,
			ClientSecret: c.ClientSecret,
		},
		s.debug,
	)
}
