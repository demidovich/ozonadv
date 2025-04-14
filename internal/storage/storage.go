package storage

import (
	"os"
	"ozonadv/pkg/utils"
)

type Storage struct {
	rootDir      string
	cabinetsFile string
	cabinets     *storageCabinets
	statsDir     string
	stats        *storageStats
}

func NewTemp() *Storage {
	return New(os.TempDir())
}

func New(rootDir string) *Storage {
	rootDir = rootDir + "/ozonadv"
	cabinetsFile := rootDir + "/cabinets.json"
	statsDir := rootDir + "/stats"

	utils.DirInit(rootDir)
	utils.DirInit(statsDir)

	return &Storage{
		rootDir:      rootDir,
		cabinetsFile: cabinetsFile,
		cabinets:     newStorageCabinets(cabinetsFile),
		statsDir:     statsDir,
		stats:        newStorageStats(statsDir),
	}
}

func (s *Storage) RootDir() string {
	return s.rootDir
}

func (s *Storage) Cabinets() *storageCabinets {
	return s.cabinets
}

func (s *Storage) Stats() *storageStats {
	return s.stats
}
