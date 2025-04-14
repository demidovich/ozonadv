package cabinets

import (
	"io"
	"ozonadv/internal/models"
	"ozonadv/internal/ozon"
	"strings"
	"time"
)

type Service struct {
	storage        storage
	campaignsCache []ozon.Campaign
}

func NewService(out io.Writer, s storage) *Service {
	return &Service{
		storage: s,
	}
}

func (s *Service) All() []models.Cabinet {
	return s.storage.All()
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
	Title string
	State string
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

	if filters.State != "" {
		filtered := []ozon.Campaign{}
		for _, campaign := range result {
			if strings.Contains(strings.ToLower(campaign.State), strings.ToLower(filters.State)) {
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
		false,
	)
}
