package storage

import (
	"cmp"
	"log"
	"ozonadv/internal/stats"
	"ozonadv/pkg/utils"
	"slices"
	"sync"
)

type storageStats struct {
	dir string
	mu  *sync.Mutex
}

func newStorageStats(dir string) *storageStats {
	s := storageStats{
		dir: dir,
		mu:  &sync.Mutex{},
	}

	return &s
}

func (s *storageStats) All() []stats.Stat {
	fnames := utils.DirListOrFail(s.dir)
	result := make([]stats.Stat, 0, len(fnames))

	for _, fname := range fnames {
		path := s.dir + "/" + fname
		stat := stats.Stat{}
		utils.JsonFileReadOrFail(path, &stat, "{}")
		if stat.UUID == "" {
			log.Fatal("некорректные данные в файле " + path)
		}
		result = append(result, stat)
	}

	slices.SortFunc(result, func(i, j stats.Stat) int {
		return cmp.Compare(i.CreatedAt, j.CreatedAt)
	})

	// sort.Slice(result, func(i, j int) bool {
	// 	return result[i].CreatedAt > result[j].CreatedAt
	// })

	return result
}

func (s *storageStats) Save(st *stats.Stat) {
	s.mu.Lock()
	defer s.mu.Unlock()

	file := s.statFile(st)
	utils.JsonFileWriteOrFail(file, st)
}

func (s *storageStats) Remove(st *stats.Stat) {
	s.mu.Lock()
	defer s.mu.Unlock()

	file := s.statFile(st)
	utils.FileRemoveOrFail(file)
}

func (s *storageStats) statFile(st *stats.Stat) string {
	return s.dir + "/" + st.UUID + ".json"
}
