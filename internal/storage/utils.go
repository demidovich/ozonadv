package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func readJsonFile(path string, result any, defaultContent string) {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	if string(content) == "" {
		content = []byte(defaultContent)
	}

	err = json.Unmarshal(content, result)
	if err != nil {
		fmt.Println(err)
	}
}

func writeJsonFile(path string, data any) {
	content, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(path, content, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func initDir(path string) {
	err := os.MkdirAll(path, 0777)
	if err != nil {
		log.Fatal(err)
	}
}

func initFile(path string) {
	f, err := os.Open(path)
	if err == nil {
		f.Close()
		return
	}

	if os.IsNotExist(err) {
		f, err = os.Create(path)
		if err != nil {
			log.Fatal(err)
		}
		f.Close()
	}
}
