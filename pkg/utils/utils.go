package utils

import (
	"encoding/json"
	"log"
	"os"
)

func JsonFileRead(path string, result any, defaultContent string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if string(content) == "" {
		content = []byte(defaultContent)
	}

	err = json.Unmarshal(content, result)
	if err != nil {
		return err
	}

	return nil
}

func JsonFileReadOrFail(path string, result any, defaultContent string) {
	err := JsonFileRead(path, result, defaultContent)
	if err != nil {
		log.Fatal(err)
	}
}

func JsonFileWrite(path string, data any) error {
	content, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, content, 0644)
	if err != nil {
		return err
	}

	return nil
}

func JsonFileWriteOrFail(path string, data any) {
	err := JsonFileWrite(path, data)
	if err != nil {
		log.Fatal(err)
	}
}

func DirInit(path string) error {
	return os.MkdirAll(path, 0777)
}

func DirInitOrFail(path string) {
	err := DirInit(path)
	if err != nil {
		log.Fatal(err)
	}
}

func DirList(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return []string{}, err
	}

	result := []string{}
	for _, e := range entries {
		if e.Name() == "." || e.Name() == ".." {
			continue
		}
		result = append(result, e.Name())
	}

	return result, nil
}

func DirListOrFail(path string) []string {
	list, err := DirList(path)
	if err != nil {
		log.Fatal(err)
	}

	return list
}

func FileInit(path string) error {
	f, err := os.Open(path)
	if err == nil {
		f.Close()
		return nil
	}

	if os.IsNotExist(err) {
		f, err = os.Create(path)
		if err != nil {
			return err
		}
		f.Close()
	}

	return err
}

func FileInitOrFail(path string) {
	err := FileInit(path)
	if err != nil {
		log.Fatal(err)
	}
}

func FileRemoveOrFail(path string) {
	err := os.Remove(path)
	if err != nil {
		log.Fatal(err)
	}

}
