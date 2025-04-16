package ui

import "errors"

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
