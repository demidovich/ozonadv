package ui

import (
	"errors"
	"fmt"

	"github.com/demidovich/ozonadv/internal/ui/colors"
)

var (
	ErrGoBack     = errors.New("go back")
	ErrFormCancel = errors.New("form cancel")
)

func isGoBack(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, ErrGoBack)
}

func isFormCanceled(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, ErrFormCancel)
}

func printError(err error) {
	fmt.Println(colors.Warning().Sprintf("%s", err.Error()))
}

func printErrorString(str string) {
	fmt.Println(colors.Warning().Sprintf("%s", str))
}
