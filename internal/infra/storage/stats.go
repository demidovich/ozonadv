package storage

import (
	"cmp"
	"log"
	"os"
	"ozonadv/internal/models"
	"ozonadv/pkg/utils"
	"slices"
	"sync"
)

type storageStats struct {
	dir          string
	downloadsDir string
	mu           *sync.Mutex
}

func newStorageStats(dir string) *storageStats {
	s := storageStats{
		dir:          dir,
		downloadsDir: dir + "/downloads",
		mu:           &sync.Mutex{},
	}

	return &s
}

func (s *storageStats) All() []models.Stat {
	fnames := utils.DirListOrFail(s.dir)
	result := make([]models.Stat, 0, len(fnames))

	for _, fname := range fnames {
		path := s.dir + "/" + fname
		stat := models.Stat{}
		utils.JsonFileReadOrFail(path, &stat, "{}")
		if stat.UUID == "" {
			log.Fatal("некорректные данные в файле " + path)
		}
		result = append(result, stat)
	}

	slices.SortFunc(result, func(i, j models.Stat) int {
		return cmp.Compare(i.CreatedAt, j.CreatedAt)
	})

	// sort.Slice(result, func(i, j int) bool {
	// 	return result[i].CreatedAt > result[j].CreatedAt
	// })

	return result
}

func (s *storageStats) Add(st *models.Stat) {
	s.mu.Lock()
	defer s.mu.Unlock()

	file := s.statFile(st)
	utils.JsonFileWriteOrFail(file, st)
}

func (s *storageStats) SaveDownloadedFile(stat *models.Stat, filename string, data []byte) {
	utils.DirInit(s.downloadsDir)
	file := s.downloadedFile(filename)

	err := os.WriteFile(file, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *storageStats) Remove(st *models.Stat) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if utils.DirExists(s.downloadsDir) {
		for _, f := range utils.DirListOrFail(s.downloadsDir) {
			file := s.downloadedFile(f)
			utils.FileRemoveOrFail(file)
		}
	}

	file := s.statFile(st)
	utils.FileRemoveOrFail(file)

	os.Remove(s.dir)
}

func (s *storageStats) statFile(st *models.Stat) string {
	return s.dir + "/" + st.UUID + ".json"
}

func (s *storageStats) downloadedFile(fname string) string {
	return s.downloadsDir + "/" + fname
}
