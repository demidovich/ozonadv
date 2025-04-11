package cabinets

type Service struct {
	storage storage
}

func NewService(s storage) *Service {
	return &Service{
		storage: s,
	}
}

func (s *Service) All() {

}

func (s *Service) Add() {

}

func (s *Service) Remove() {

}
