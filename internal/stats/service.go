package stats

import (
	"io"
	"ozonadv/internal/models"
	"ozonadv/internal/ozon"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	out     io.Writer
	storage storage
}

type CreateOptions struct {
}

func NewService(s storage) *Service {
	return &Service{
		storage: s,
	}
}

func (s *Service) Create(options models.StatOptions, campaigns []ozon.Campaign) (*models.Stat, error) {
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
	downloader := newDownloader(s.out, st, ozonApi, s.storage)
	downloader.Start()
}

func (s *Service) ozonApi(st *models.Stat) *ozon.Ozon {
	return ozon.New(
		ozon.Config{
			ClientId:     st.Options.CabinetClientId,
			ClientSecret: st.Options.CabinetClientSecret,
		},
		false,
	)
}
