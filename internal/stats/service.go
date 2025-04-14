package stats

import (
	"io"
	"ozonadv/internal/models"
	"ozonadv/internal/ozon"
)

type Service struct {
	out     io.Writer
	storage storage
}

func NewService(out io.Writer, s storage) *Service {
	return &Service{
		out:     out,
		storage: s,
	}
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
