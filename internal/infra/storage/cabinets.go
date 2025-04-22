package storage

import (
	"cmp"
	"ozonadv/internal/models"
	"ozonadv/pkg/utils"
	"slices"
	"time"
)

type storageCabinets struct {
	file string
	data map[string]models.Cabinet
}

func newStorageCabinets(file string) *storageCabinets {
	s := storageCabinets{
		file: file,
	}

	utils.FileInitOrFail(file)
	utils.JSONFileReadOrFail(file, &s.data, "{}")

	return &s
}

func (s *storageCabinets) All() []models.Cabinet {
	result := make([]models.Cabinet, 0, len(s.data))
	for _, c := range s.data {
		result = append(result, c)
	}

	slices.SortFunc(result, func(i, j models.Cabinet) int {
		return cmp.Compare(i.CreatedAt, j.CreatedAt)
	})

	return result
}

func (s *storageCabinets) Has(cabinet models.Cabinet) bool {
	_, ok := s.data[cabinet.UUID]
	return ok
}

func (s *storageCabinets) Add(cabinet models.Cabinet) {
	cabinet.CreatedAt = time.Now().String()
	s.data[cabinet.UUID] = cabinet

	s.saveStorage()
}

func (s *storageCabinets) Remove(cabinet models.Cabinet) {
	delete(s.data, cabinet.UUID)

	s.saveStorage()
}

func (s *storageCabinets) saveStorage() {
	utils.JSONFileWriteOrFail(s.file, s.data)
}
