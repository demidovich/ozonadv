package config

import (
	"errors"
	"fmt"
	"os"
	"ozonstat/internal/google"
	"ozonstat/internal/ozon"

	"github.com/spf13/viper"
)

type Config struct {
	Ozon   ozon.Config
	Google google.Config
}

func NewOrFail(filename string) Config {
	cfg, err := New(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return cfg
}

func New(filename string) (Config, error) {
	v := viper.New()
	v.SetConfigFile(filename)
	v.AddConfigPath("..")
	v.AutomaticEnv()

	cfg := Config{}

	if err := v.ReadInConfig(); err != nil {
		return cfg, err
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return cfg, fmt.Errorf("unable to decode config into struct, %v", err)
	}

	if cfg.Ozon.ClientId == "" {
		return cfg, errors.New("В config.yml отсутствует или пустой параметр ozon.ClientId")
	}

	if cfg.Ozon.ClientSecret == "" {
		return cfg, errors.New("В config.yml отсутствует или пустой параметр ozon.ClientSecret")
	}

	return cfg, nil
}
