package stats

import "ozonadv/internal/models"

type storage interface {
	All() []models.Stat

	Add(stat *models.Stat)

	Remove(stat *models.Stat)

	SaveDownloadedFile(stat *models.Stat, filename string, data []byte)
}
