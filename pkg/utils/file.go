package utils

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

func JSONFileRead(path string, result any, defaultContent string) error {
	path = filepath.Clean(path)
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if len(content) == 0 {
		content = []byte(defaultContent)
	}

	err = json.Unmarshal(content, result)
	if err != nil {
		return err
	}

	return nil
}

func JSONFileReadOrFail(path string, result any, defaultContent string) {
	err := JSONFileRead(path, result, defaultContent)
	if err != nil {
		log.Fatal(err)
	}
}

func JSONFileWrite(path string, data any) error {
	path = filepath.Clean(path)
	content, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, content, 0600) // 0644
	if err != nil {
		return err
	}

	return nil
}

func JSONFileWriteOrFail(path string, data any) {
	err := JSONFileWrite(path, data)
	if err != nil {
		log.Fatal(err)
	}
}

func DirInit(path string) error {
	path = filepath.Clean(path)
	return os.MkdirAll(path, 0750) // 0777
}

func DirInitOrFail(path string) {
	err := DirInit(path)
	if err != nil {
		log.Fatal(err)
	}
}

func DirExists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

func DirList(path string) ([]string, error) {
	path = filepath.Clean(path)
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
	path = filepath.Clean(path)

	f, err := os.Open(path)
	if err == nil {
		return f.Close()
	}

	if os.IsNotExist(err) {
		f, err = os.Create(path)
		if err == nil {
			return f.Close()
		}
		return err
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
	path = filepath.Clean(path)
	err := os.Remove(path)
	if err != nil {
		log.Fatal(err)
	}
}
