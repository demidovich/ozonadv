package storage

import (
	"log"
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
