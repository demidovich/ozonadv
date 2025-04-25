package stats

import "github.com/demidovich/ozonadv/internal/models"

type Storage interface {
	All() []*models.Stat

	Add(stat *models.Stat)

	Remove(stat *models.Stat)

	AddDownloadsFile(stat *models.Stat, filename string, data []byte)

	ReadDownloadedFile(stat *models.Stat, filename string) []byte
}
