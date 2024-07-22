package mediaconv

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port int `json:"port"`
}

func LoadConfig(filename string) (*Config, error) {
	var config *Config
	configBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(configBytes, &config); err != nil {
		return nil, err
	}

	return config, nil
}
