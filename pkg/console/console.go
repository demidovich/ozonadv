package console

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
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

func InputString(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt + ": ")
	text, _ := reader.ReadString('\n')

	return strings.TrimRight(text, "\n")
}

func InputInt(prompt string) int {
	s := InputString(prompt)
	i, _ := strconv.Atoi(s)

	return i
}
