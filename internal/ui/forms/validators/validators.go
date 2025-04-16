package validators

import (
	"errors"
	"regexp"
)

var (
	regexpDate *regexp.Regexp
)

func init() {
	regexpDate = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
}

func Required(s string) error {
	if s == "" {
		return errors.New("поле не заполнено")
	}

	return nil
}

func Date(s string) error {
	if !regexpDate.MatchString(s) {
		return errors.New("поле должно быть в формате ГГГГ-ДД-ММ")
	}

	return nil
}

func DateRequiured(s string) error {
	if err := Required(s); err != nil {
		return err
	}

	if !regexpDate.MatchString(s) {
		return errors.New("поле должно быть в формате ГГГГ-ДД-ММ")
	}

	return nil
}
