package stats

import (
	"ozonadv/internal/infra/ozon"
	"ozonadv/internal/models"
	"time"

	"github.com/google/uuid"
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

func (s *Service) All() []models.Stat {
	return s.storage.All()
}

func (s *Service) CabinetAll(cabinet models.Cabinet) []models.Stat {
	result := []models.Stat{}
	for _, s := range s.storage.All() {
		if cabinet.UUID == s.Options.CabinetUUID {
			result = append(result, s)
		}
	}

	return result
}

func (s *Service) Find(uuid string) (*models.Stat, bool) {
	for _, c := range s.storage.All() {
		if c.UUID == uuid {
			return &c, true
		}
	}
	return nil, false
}

func (s *Service) Create(options models.StatOptions, campaigns []models.Campaign) (*models.Stat, error) {
	st := &models.Stat{}
	if err := options.Validate(); err != nil {
		return st, err
	}

	st.UUID = uuid.New().String()
	st.Options = options
	st.CreatedAt = time.Now().String()

	for _, c := range campaigns {
		st.AddCampaign(c)
	}

	s.storage.Add(st)

	return st, nil
}

func (s *Service) Download(st *models.Stat) {
	ozonApi := s.ozonApi(st)
	downloader := newDownloader(st, ozonApi, s.storage, s.debug)
	downloader.Start()
}

func (s *Service) ExportToFile(stat *models.Stat, file string) error {
	e := newExport(s.storage, s.debug)

	return e.toFile(stat, file)
}

func (s *Service) Remove(stat *models.Stat) {
	s.storage.Remove(stat)
}

func (s *Service) ozonApi(st *models.Stat) *ozon.Ozon {
	return ozon.New(
		ozon.Config{
			ClientId:     st.Options.CabinetClientId,
			ClientSecret: st.Options.CabinetClientSecret,
		},
		s.debug,
	)
}
