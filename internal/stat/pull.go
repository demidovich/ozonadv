package stat

import (
	"fmt"
	"ozonadv/internal/ozon"
	"ozonadv/internal/storage"

	"github.com/go-playground/validator/v10"
)

type PullOptions struct {
	ExportFile string `validate:"required,filepath"`
}

func (c *PullOptions) validate() error {
	validate := validator.New()

	err := validate.Struct(c)
	if err == nil {
		return nil
	}

	errs := err.(validator.ValidationErrors)
	if errs != nil {
		return fmt.Errorf("%s", errs)
	}

	return nil
}

type pullUsecase struct {
	storage    *storage.Storage
	ozonClient *ozon.Client
}

func (p *pullUsecase) Handle(options PullOptions) error {
	return nil
}
