package storage

import (
	"cmp"
	"log"
	"os"
	"path/filepath"
	"slices"

	"github.com/demidovich/ozonadv/internal/models"
	"github.com/demidovich/ozonadv/pkg/utils"
)

type statsStorage struct {
	dir          string
	downloadsDir string
}

func newStatsStorage(dir string) *statsStorage {
	s := statsStorage{
		dir:          dir,
		downloadsDir: dir + "/downloads",
	}

	return &s
}

func (s *statsStorage) All() []*models.Stat {
	fnames := utils.DirListOrFail(s.dir)
	result := make([]*models.Stat, 0, len(fnames))

	for _, fname := range fnames {
		if fname == "downloads" {
			continue
		}

		path := s.dir + "/" + fname
		stat := models.Stat{}
		utils.JSONFileReadOrFail(path, &stat, "{}")
		if stat.UUID == "" {
			log.Fatal("некорректные данные в файле " + path)
		}
		result = append(result, &stat)
	}

	slices.SortFunc(result, func(i, j *models.Stat) int {
		return cmp.Compare(i.CreatedAt, j.CreatedAt)
	})

	// sort.Slice(result, func(i, j int) bool {
	// 	return result[i].CreatedAt > result[j].CreatedAt
	// })

	return result
}

func (s *statsStorage) Add(st *models.Stat) {
	file := s.statFile(st)
	utils.JSONFileWriteOrFail(file, st)
}

func (s *statsStorage) AddDownloadsFile(stat *models.Stat, filename string, data []byte) {
	utils.DirInitOrFail(s.downloadsDir)
	file := s.downloadedFile(filename)

	err := os.WriteFile(file, data, 0600)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *statsStorage) ReadDownloadedFile(stat *models.Stat, filename string) []byte {
	file := s.downloadedFile(filename)
	file = filepath.Clean(file)

	content, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	return content
}

func (s *statsStorage) Remove(st *models.Stat) {
	if utils.DirExists(s.downloadsDir) {
		for _, f := range utils.DirListOrFail(s.downloadsDir) {
			file := s.downloadedFile(f)
			utils.FileRemoveOrFail(file)
		}
	}

	file := s.statFile(st)
	utils.FileRemoveOrFail(file)

	if err := os.Remove(s.dir); err != nil {
		log.Fatal(err)
	}
}

func (s *statsStorage) statFile(st *models.Stat) string {
	return s.dir + "/" + st.UUID + ".json"
}

func (s *statsStorage) downloadedFile(fname string) string {
	return s.downloadsDir + "/" + fname
}
