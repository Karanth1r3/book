package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type (
	DB struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		UserName string `yaml:"username"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	}
	Config struct {
		DB DB `yaml:"db"`
	}
)

func Parse(file string) (*Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %w", err)
	}
	cfg := Config{}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal data from file to cfg: %w", err)
	}
	return &cfg, nil
}
