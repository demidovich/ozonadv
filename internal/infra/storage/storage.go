package storage

import (
	"log"
	"os"

	"github.com/demidovich/ozonadv/pkg/utils"
)

type Storage struct {
	rootDir      string
	cabinetsFile string
	cabinets     *cabinetsStorage
	statsDir     string
	stats        *statsStorage
}

func NewDefault() *Storage {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return New(homedir)
}

func New(rootDir string) *Storage {
	rootDir += "/.ozonadv"
	cabinetsFile := rootDir + "/cabinets.json"
	statsDir := rootDir + "/stats"

	utils.DirInitOrFail(rootDir)
	utils.DirInitOrFail(statsDir)

	return &Storage{
		rootDir:      rootDir,
		cabinetsFile: cabinetsFile,
		cabinets:     newCabinetsStorage(cabinetsFile),
		statsDir:     statsDir,
		stats:        newStatsStorage(statsDir),
	}
}

func (s *Storage) RootDir() string {
	return s.rootDir
}

func (s *Storage) Cabinets() *cabinetsStorage {
	return s.cabinets
}

func (s *Storage) Stats() *statsStorage {
	return s.stats
}
