package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type GymPortalConfig struct {
	URL     string
	Cookies string
}

type InfluxDBConfig struct {
	URL       string
	AuthToken string
	Org       string
	Bucket    string
}

type Config struct {
	GymPortal GymPortalConfig
	InfluxDB  InfluxDBConfig
}

func Load(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("couldn't read the file %s: %w", path, err)
	}

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return Config{}, fmt.Errorf("invalid json: %w", err)
	}

	return config, nil
}
