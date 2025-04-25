package storage

import (
	"cmp"
	"slices"
	"time"

	"github.com/demidovich/ozonadv/internal/models"
	"github.com/demidovich/ozonadv/pkg/utils"
)

type cabinetsStorage struct {
	file string
	data map[string]models.Cabinet
}

func newCabinetsStorage(file string) *cabinetsStorage {
	s := cabinetsStorage{
		file: file,
	}

	utils.FileInitOrFail(file)
	utils.JSONFileReadOrFail(file, &s.data, "{}")

	return &s
}

func (s *cabinetsStorage) All() []models.Cabinet {
	result := make([]models.Cabinet, 0, len(s.data))
	for _, c := range s.data {
		result = append(result, c)
	}

	slices.SortFunc(result, func(i, j models.Cabinet) int {
		return cmp.Compare(i.CreatedAt, j.CreatedAt)
	})

	return result
}

func (s *cabinetsStorage) Has(cabinet models.Cabinet) bool {
	_, ok := s.data[cabinet.UUID]
	return ok
}

func (s *cabinetsStorage) Add(cabinet models.Cabinet) {
	cabinet.CreatedAt = time.Now().String()
	s.data[cabinet.UUID] = cabinet

	s.saveStorage()
}

func (s *cabinetsStorage) Remove(cabinet models.Cabinet) {
	delete(s.data, cabinet.UUID)

	s.saveStorage()
}

func (s *cabinetsStorage) saveStorage() {
	utils.JSONFileWriteOrFail(s.file, s.data)
}
