package cabinets

import "github.com/demidovich/ozonadv/internal/models"

type storage interface {
	All() []models.Cabinet

	Has(cabinet models.Cabinet) bool

	Add(cabinet models.Cabinet)

	Remove(cabinet models.Cabinet)
}
