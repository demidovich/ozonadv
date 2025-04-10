package storage

import (
	"cmp"
	"ozonadv/internal/cabinets"
	"ozonadv/pkg/utils"
	"slices"
	"time"
)

type storageCabinets struct {
	file string
	data map[string]cabinets.Cabinet
}

func newStorageCabinets(file string) *storageCabinets {
	s := storageCabinets{
		file: file,
	}

	utils.FileInitOrFail(file)
	utils.JsonFileReadOrFail(file, &s.data, "{}")

	return &s
}

func (s *storageCabinets) All() []cabinets.Cabinet {
	result := make([]cabinets.Cabinet, 0, len(s.data))
	for _, c := range s.data {
		result = append(result, c)
	}

	slices.SortFunc(result, func(i, j cabinets.Cabinet) int {
		return cmp.Compare(i.CreatedAt, j.CreatedAt)
	})

	return result
}

func (s *storageCabinets) Save(cabinet cabinets.Cabinet) {
	cabinet.CreatedAt = time.Now().String()
	s.data[cabinet.UUID] = cabinet

	s.saveStorage()
}

func (s *storageCabinets) Remove(cabinet cabinets.Cabinet) {
	delete(s.data, cabinet.UUID)

	s.saveStorage()
}

func (s *storageCabinets) saveStorage() {
	utils.JsonFileWriteOrFail(s.file, s.data)
}
