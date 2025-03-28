package stat_export

import (
	"ozonadv/internal/storage"
)

type statExport struct {
	storage *storage.Storage
}

func New(storage *storage.Storage) statExport {
	return statExport{
		storage: storage,
	}
}

func (s statExport) ToFile(file string) error {
	return nil
}
