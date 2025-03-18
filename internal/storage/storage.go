package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"ozonadv/internal/ozon"
)

type Storage struct {
	rootDir        string
	statisticsFile string
}

func New() *Storage {
	fmt.Println("Инициализация локального хранилища")

	root := os.TempDir() + "/ozonadv"
	s := Storage{
		rootDir:        root,
		statisticsFile: root + "/statistics.json",
	}

	initDir(s.rootDir)
	initFile(s.statisticsFile)

	fmt.Println("Директория локального хранилища", s.rootDir)
	fmt.Println("")

	return &s
}

func (s *Storage) SaveStatistic(item ozon.Statistic) {
	all := *s.FindAllStatistic()

	updated := false
	for i, row := range all {
		if row.UUID == item.UUID {
			all[i] = item
			updated = true
		}
	}

	if !updated {
		all = append(all, item)
	}

	writeJsonFile(s.statisticsFile, &all)
}

func (s *Storage) FindStatistic(uuid string) (*ozon.Statistic, error) {
	all := s.FindAllStatistic()

	for _, row := range *all {
		if row.UUID == uuid {
			return &row, nil
		}
	}

	return nil, errors.New("not found")
}

func (s *Storage) FindAllStatistic() *[]ozon.Statistic {
	result := []ozon.Statistic{}
	readJsonFile(s.statisticsFile, &result)

	return &result
}

func readJsonFile(path string, result any) {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
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
