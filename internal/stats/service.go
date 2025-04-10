package stats

import (
	"io"
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

func (s *Service) Download(st *Stat) {
	ozonApi := s.ozonApi(st)
	downloader := newDownloader(s.out, st, ozonApi, s.storage)
	downloader.Start()
}

func (s *Service) ozonApi(st *Stat) *ozon.Ozon {
	return ozon.New(
		ozon.Config{
			ClientId:     st.Options.CabinetClientId,
			ClientSecret: st.Options.CabinetClientSecret,
		},
		false,
	)
}
