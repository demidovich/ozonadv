package cabinets

import (
	"io"
	"log"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/demidovich/ozonadv/internal/infra/ozon"
	"github.com/demidovich/ozonadv/internal/models"
)

type Service struct {
	storage        storage
	campaignsCache map[string][]models.Campaign
	debug          Debug
}

func NewService(out io.Writer, storage storage, debug Debug) *Service {
	return &Service{
		storage:        storage,
		campaignsCache: make(map[string][]models.Campaign),
		debug:          debug,
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

func (f *CampaignFilters) ids() map[string]bool {
	v := map[string]bool{}
	matched, err := regexp.MatchString(`^\s*\d+\s*(\s*,\s*\d+\s*)*$`, f.Title)
	if err != nil {
		log.Fatal(err)
	}

	if !matched {
		return v
	}

	s := strings.ReplaceAll(f.Title, " ", "")
	s = strings.Trim(s, ",")

	for _, id := range strings.Split(s, ",") {
		v[id] = true
	}

	return v
}

func (s *Service) Campaigns(cabinet models.Cabinet) ([]models.Campaign, error) {
	if cache, ok := s.campaignsCache[cabinet.UUID]; ok {
		return cache, nil
	}

	result, err := s.ozon(cabinet).Campaigns().All()
	if err != nil {
		return result, err
	}

	s.campaignsCache[cabinet.UUID] = result

	return result, nil
}

func (s *Service) CampaignsFiltered(cabinet models.Cabinet, filters CampaignFilters) ([]models.Campaign, error) {
	result, err := s.ozon(cabinet).Campaigns().All()
	if err != nil {
		return result, err
	}

	ids := filters.ids()
	if len(ids) > 0 {
		filtered := []models.Campaign{}
		for _, campaign := range result {
			if _, ok := ids[campaign.ID]; ok {
				filtered = append(filtered, campaign)
			}
		}
		result = filtered
	} else if filters.Title != "" {
		filtered := []models.Campaign{}
		for _, campaign := range result {
			if strings.Contains(strings.ToLower(campaign.Title), strings.ToLower(filters.Title)) {
				filtered = append(filtered, campaign)
			}
		}
		result = filtered
	}

	if len(filters.States) > 0 {
		filtered := []models.Campaign{}
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
			ClientID:     c.ClientID,
			ClientSecret: c.ClientSecret,
		},
		s.debug,
	)
}
