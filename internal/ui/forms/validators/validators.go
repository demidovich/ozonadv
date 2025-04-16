package validators

import "errors"

func Required(s string) error {
	if s == "" {
		return errors.New("поле не заполнено")
	}
	return nil
}
