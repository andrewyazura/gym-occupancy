package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLoadConfigNoFile(t *testing.T) {
	config_path := fmt.Sprintf("%s/config.temp.json", t.TempDir())

	_, err := loadConfig(config_path)

	err = errors.Unwrap(err)
	if err, ok := err.(*os.PathError); !ok {
		t.Errorf("invalid error, expected *os.PathError, received: %v", err)
	}
}

func TestLoadConfigInvalidJSON(t *testing.T) {
	config_path := fmt.Sprintf("%s/config.temp.json", t.TempDir())

	file, err := os.OpenFile(config_path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		t.Fatalf("couldn't create a test config file: %v", err)
	}

	if _, err := file.WriteString(`aaaaa`); err != nil {
		t.Fatalf("couldn't write to file: %v", err)
	}

	if err := file.Close(); err != nil {
		t.Fatalf("couldn't close config file: %T %v", err, err)
	}

	_, err = loadConfig(config_path)

	err = errors.Unwrap(err)
	if err, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("invalid error, expected *json.SyntaxError, got: %T, %v", err, err)
	}
}

func TestLoadConfig(t *testing.T) {
	expected_config := Config{
		URL:        "https://example.com",
		Cookies:    "ClientPortal.Auth=auth-token",
		Timezone:   "Europe/London",
		Keywords:   []string{"wrocławska", "korfantego"},
		OutputFile: "results.csv",
	}

	content := `{
  "url": "https://example.com",
  "cookies": "ClientPortal.Auth=auth-token",
  "timezone": "Europe/London",
  "keywords": ["wrocławska", "korfantego"],
  "outputFile": "results.csv"
}`

	config_path := fmt.Sprintf("%s/config.temp.json", t.TempDir())

	file, err := os.OpenFile(config_path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		t.Fatalf("couldn't create a test config file: %v", err)
	}

	if _, err := file.WriteString(content); err != nil {
		t.Fatalf("couldn't write to file: %v", err)
	}

	if err := file.Close(); err != nil {
		t.Fatalf("couldn't close config file: %v", err)
	}

	config, err := loadConfig(config_path)
	if err != nil {
		t.Fatalf("couldn't load config: %v", err)
	}

	if config.URL != expected_config.URL {
		t.Errorf("invalid URL, expected %s, got %s", expected_config.URL, config.URL)
	}

	if config.Cookies != expected_config.Cookies {
		t.Errorf("invalid cookies, expected %s, got %s", expected_config.Cookies, config.Cookies)
	}

	if config.Timezone != expected_config.Timezone {
		t.Errorf("invalid timezone, expected %s, got %s", expected_config.Timezone, config.Timezone)
	}

	if !cmp.Equal(config.Keywords, expected_config.Keywords) {
		t.Errorf("invalid keywords, expected %v, got %v", expected_config.Keywords, config.Keywords)
	}

	if config.OutputFile != expected_config.OutputFile {
		t.Errorf("invalid output file, expected %s, got %s", expected_config.OutputFile, config.OutputFile)
	}
}

func TestAppendToOutputNoFile(t *testing.T) {
	output_path := fmt.Sprintf("%s/output.temp.csv", t.TempDir())

	row := []byte("test\n")
	err := appendToOutput(output_path, row)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	file, err := os.Open(output_path)
	if err != nil {
		t.Fatalf("couldn't open output file: %v", err)
	}

	stat, err := file.Stat()
	if err != nil {
		t.Fatalf("couldn't get stats of output file: %v", err)
	}

	if stat.Size() == 0 {
		t.Errorf("output file is empty")
	}
}

func TestAppendToOutput(t *testing.T) {
	output_path := fmt.Sprintf("%s/output.temp.csv", t.TempDir())

	row := []byte("test\n")

	if err := appendToOutput(output_path, row); err != nil {
		t.Fatalf("couldn't append to output file: %v", err)
	}

	file, err := os.Open(output_path)
	if err != nil {
		t.Fatalf("couldn't open output file: %v", err)
	}

	stat, err := file.Stat()
	if err != nil {
		t.Fatalf("couldn't get stats of output file: %v", err)
	}

	if stat.Size() == 0 {
		t.Errorf("output file is empty")
	}
}
