package console

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func Ask(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [Да/Нет]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "д" || response == "да" || response == "y" || response == "yes" {
			return true
		}

		if response == "н" || response == "нет" || response == "n" || response == "no" {
			return false
		}
	}
}
