package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir += "/gatorconfig.json"
	return dir, nil
}

func Read() (*Config, error) {
	cfg := Config{}

	dir, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(dir)
	if err != nil {
		return nil, err
	}

	e := json.Unmarshal(data, &cfg)
	if e != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) SetUser(name string) error {
	c.CurrentUserName = name
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	e := os.WriteFile(path, data, 0644)
	if e != nil {
		return e
	}
	return nil
}
