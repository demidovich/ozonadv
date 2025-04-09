package config

import (
	"errors"
	"fmt"
	"log"
	"ozonadv/internal/google"
	"ozonadv/internal/ozon"

	"github.com/spf13/viper"
)

type Config struct {
	Verbose bool
	Ozon    ozon.Config
	Google  google.Config
}

func NewOrFail(filename string) Config {
	cfg, err := New(filename)
	if err != nil {
		fmt.Println("")
		fmt.Println("Ошибка конфигурации: не найден файл config.yaml")
		fmt.Println("В директории с генератором статистики ozonadv должен находиться корректный yaml-файл config.yaml")
		fmt.Println("")
		fmt.Println("Описание конфигурации и пример файла можно посмотреть на github проекта")
		fmt.Println("https://github.com/demidovich/ozonadv")
		fmt.Println("")
		log.Fatal(err)
	}

	return cfg
}

func New(filename string) (Config, error) {
	v := viper.New()
	v.SetConfigFile(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()

	cfg := Config{}

	if err := v.ReadInConfig(); err != nil {
		return cfg, err
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return cfg, fmt.Errorf("ошибка парсинга конфигурационного файла, %v", err)
	}

	if cfg.Ozon.ClientId == "" {
		return cfg, errors.New("в config.yml отсутствует или пустой параметр ozon.ClientId")
	}

	if cfg.Ozon.ClientSecret == "" {
		return cfg, errors.New("в config.yml отсутствует или пустой параметр ozon.ClientSecret")
	}

	return cfg, nil
}
