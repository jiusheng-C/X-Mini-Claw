package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	APIKey                string `json:"api_key"`
	BaseURL               string `json:"base_url"`
	Model                 string `json:"model"`
	APIStyle              string `json:"api_style"`
	Workspace             string `json:"workspace"`
	CommandTimeoutSeconds int    `json:"command_timeout_seconds"`
}

func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
