package stat

import (
	"fmt"
	"ozonstat/internal/ozon"
)

type FetchOptions struct {
	Days uint
}

func Fetch(ozonClient *ozon.Client, options FetchOptions) error {

	fmt.Println(options)
	fmt.Println(ozonClient)

	return nil
}
