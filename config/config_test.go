package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"
)

func TestLoadConfigNoFile(t *testing.T) {
	configPath := fmt.Sprintf("%s/config.json", t.TempDir())

	_, err := Load(configPath)

	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("invalid error, expected os.ErrNotExist, received: %v", err)
	}
}

func TestLoadConfigInvalidJSON(t *testing.T) {
	configPath := fmt.Sprintf("%s/config.json", t.TempDir())

	file, err := os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		t.Fatalf("couldn't create a test config file: %v", err)
	}

	if _, err := file.WriteString(`aaaaa`); err != nil {
		t.Fatalf("couldn't write to file: %v", err)
	}

	if err := file.Close(); err != nil {
		t.Fatalf("couldn't close config file: %T %v", err, err)
	}

	_, err = Load(configPath)

	var syntaxErr *json.SyntaxError
	if !errors.As(err, &syntaxErr) {
		t.Errorf("invalid error, expected *json.SyntaxError, got: %T, %v", err, err)
	}
}

func TestLoadConfig(t *testing.T) {
	expectedConfig := Config{
		GymPortal: GymPortalConfig{
			URL:     "https://example.com",
			Cookies: "ClientPortal.Auth=auth-token",
		},
		InfluxDB: InfluxDBConfig{
			URL:       "http://localhost:8086",
			AuthToken: "auth-token",
			Org:       "org",
			Bucket:    "bucket",
		},
	}

	content := `{
    "gymPortal": {
      "url": "https://example.com",
      "cookies": "ClientPortal.Auth=auth-token"
    },
    "influxDB": {
      "url": "http://localhost:8086",
      "authToken": "auth-token",
      "org": "org",
      "bucket": "bucket"
    }
  }`

	configPath := fmt.Sprintf("%s/config.json", t.TempDir())

	file, err := os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		t.Fatalf("couldn't create a test config file: %v", err)
	}

	if _, err := file.WriteString(content); err != nil {
		t.Fatalf("couldn't write to file: %v", err)
	}

	if err := file.Close(); err != nil {
		t.Fatalf("couldn't close config file: %v", err)
	}

	config, err := Load(configPath)
	if err != nil {
		t.Fatalf("couldn't load config: %v", err)
	}

	fmt.Printf("%v\n", config)
	if config != expectedConfig {
		t.Fatalf("error")
	}
}
