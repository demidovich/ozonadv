package validation

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ru_translations "github.com/go-playground/validator/v10/translations/ru"
)

var validatorInstance *validator.Validate
var translatorInstance ut.Translator

func init() {
	validatorInstance = validator.New()

	rus := ru.New()
	uni := ut.New(rus, rus)

	translatorInstance, _ = uni.GetTranslator("ru")
	err := ru_translations.RegisterDefaultTranslations(validatorInstance, translatorInstance)
	if err != nil {
		log.Fatal(err)
	}
}

func ValidateStruct(s any) error {
	err := validatorInstance.Struct(s)
	if err == nil {
		return nil
	}

	errs := err.(validator.ValidationErrors)
	if errs != nil {
		m := []string{"", ""}
		for _, e := range errs {
			m = append(m, e.Translate(translatorInstance))
		}
		m = append(m, "")

		return fmt.Errorf("%s", strings.Join(m, "\n"))
	}

	return nil
}
