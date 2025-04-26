package factory

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/demidovich/ozonadv/internal/models"
)

type cabinetFactory struct {
}

func Cabinet() *cabinetFactory {
	return &cabinetFactory{}
}

func (f *cabinetFactory) New() models.Cabinet {
	c := models.Cabinet{
		UUID:         gofakeit.UUID(),
		Name:         gofakeit.Company(),
		ClientID:     gofakeit.UUID(),
		ClientSecret: gofakeit.UUID(),
		CreatedAt:    time.Now().String(),
	}

	return c
}
