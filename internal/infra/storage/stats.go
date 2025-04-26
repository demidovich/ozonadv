package storage

import (
	"cmp"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"slices"

	"github.com/demidovich/ozonadv/internal/models"
	"github.com/demidovich/ozonadv/pkg/utils"
)

type statsStorage struct {
	dir string
}

func newStatsStorage(dir string) *statsStorage {
	s := statsStorage{
		dir: dir,
	}

	return &s
}

func (s *statsStorage) All() []*models.Stat {
	fnames := utils.DirListOrFail(s.dir)
	result := make([]*models.Stat, 0, len(fnames))

	for _, statUUID := range fnames {
		file := s.statFile(statUUID)
		stat := models.Stat{}

		utils.JSONFileReadOrFail(file, &stat, "{}")
		if stat.UUID == "" {
			log.Fatal("некорректные данные в файле " + file)
		}

		result = append(result, &stat)
	}

	slices.SortFunc(result, func(i, j *models.Stat) int {
		return cmp.Compare(i.CreatedAt, j.CreatedAt)
	})

	return result
}

func (s *statsStorage) Add(stat *models.Stat) {
	dir := s.statDir(stat.UUID)
	utils.DirInitOrFail(dir)

	file := s.statFile(stat.UUID)
	utils.JSONFileWriteOrFail(file, stat)
}

func (s *statsStorage) AddDownloadsFile(stat *models.Stat, filename string, data []byte) {
	dir := s.downloadsDir(stat.UUID)
	utils.DirInitOrFail(dir)

	file := s.downloadsFile(stat.UUID, filename)

	err := os.WriteFile(file, data, 0600)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *statsStorage) ReadDownloadsFile(stat *models.Stat, filename string) []byte {
	file := s.downloadsFile(stat.UUID, filename)
	file = filepath.Clean(file)

	content, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	return content
}

func (s *statsStorage) Remove(stat *models.Stat) {
	statDir := s.statDir(stat.UUID)
	statDir = path.Clean(statDir)

	os.RemoveAll(statDir)
}

// Generate directory file path
func (s *statsStorage) statDir(statUUID string) string {
	return fmt.Sprintf("%s/%s", s.dir, statUUID)
}

// Generate stat.json file path
func (s *statsStorage) statFile(statUUID string) string {
	return fmt.Sprintf("%s/%s/stat.json", s.dir, statUUID)
}

// Generate downloads directory path
func (s *statsStorage) downloadsDir(statUUID string) string {
	return fmt.Sprintf("%s/%s/downloads", s.dir, statUUID)
}

// Generate downloads file path
func (s *statsStorage) downloadsFile(statUUID string, fname string) string {
	return fmt.Sprintf("%s/%s/downloads/%s", s.dir, statUUID, fname)
}
