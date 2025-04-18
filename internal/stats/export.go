package stats

import (
	"errors"
	"ozonadv/internal/models"
)

type export struct {
	storage Storage
	debug   Debug
}

func newExport(storage Storage, debug Debug) export {
	return export{
		storage: storage,
		debug:   debug,
	}
}

func (e export) toFile(stat *models.Stat, file string) error {
	if stat.IsTypeObject() {
		return newObjectStatExport(e.storage, e.debug, stat).toFile(file)
	}

	if stat.IsTypeTotal() {
		return errors.New("отчет total еще не реализован")
	}

	return nil
}
