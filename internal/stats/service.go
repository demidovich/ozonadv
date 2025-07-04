package stats

import (
	"time"

	"github.com/demidovich/ozonadv/internal/infra/ozon"
	"github.com/demidovich/ozonadv/internal/models"

	googleUUID "github.com/google/uuid"
)

type Service struct {
	debug   Debug
	storage Storage
}

type CreateOptions struct {
}

func NewService(s Storage, d Debug) *Service {
	return &Service{
		debug:   d,
		storage: s,
	}
}

func (s *Service) All() []*models.Stat {
	return s.storage.All()
}

func (s *Service) CabinetAll(cabinet models.Cabinet) []*models.Stat {
	result := []*models.Stat{}
	for _, s := range s.storage.All() {
		if cabinet.UUID == s.Options.CabinetUUID {
			result = append(result, s)
		}
	}

	return result
}

func (s *Service) Find(uuid string) (*models.Stat, bool) {
	for _, stat := range s.storage.All() {
		if stat.UUID == uuid {
			return stat, true
		}
	}
	return nil, false
}

func (s *Service) FindByRequestUUID(requestUUID string) (*models.Stat, bool) {
	for _, stat := range s.storage.All() {
		for _, item := range stat.Items {
			if item.Request.UUID == requestUUID {
				return stat, true
			}
		}
	}

	return nil, false
}

func (s *Service) Create(options models.StatOptions, campaigns []models.Campaign) (*models.Stat, error) {
	st := &models.Stat{}
	if err := options.Validate(); err != nil {
		return st, err
	}

	st.UUID = googleUUID.New().String()
	st.Options = options
	st.CreatedAt = time.Now().String()

	for _, c := range campaigns {
		st.AddCampaign(c)
	}

	s.storage.Add(st)

	return st, nil
}

func (s *Service) Download(st *models.Stat) {
	ozonAPI := s.ozonAPI(st)
	downloader := newDownloader(st, ozonAPI, s.storage, s.debug)
	downloader.Start()
}

func (s *Service) ExportToFile(stat *models.Stat, file string) error {
	e := newExport(s.storage, s.debug)

	return e.toFile(stat, file)
}

func (s *Service) Remove(stat *models.Stat) {
	s.storage.Remove(stat)
}

func (s *Service) ozonAPI(st *models.Stat) *ozon.Ozon {
	return ozon.New(
		ozon.Config{
			ClientID:     st.Options.CabinetClientID,
			ClientSecret: st.Options.CabinetClientSecret,
		},
		s.debug,
	)
}
