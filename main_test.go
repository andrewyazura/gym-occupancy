package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLoadConfigNoFile(t *testing.T) {
	configPath := fmt.Sprintf("%s/config.json", t.TempDir())

	_, err := loadConfig(configPath)

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

	_, err = loadConfig(configPath)

	var syntaxErr *json.SyntaxError
	if !errors.As(err, &syntaxErr) {
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

	config, err := loadConfig(configPath)
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
	outputPath := fmt.Sprintf("%s/output.csv", t.TempDir())

	row := []byte("test\n")
	err := appendToOutput(outputPath, row)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	file, err := os.Open(outputPath)
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

func TestAppendToOutputReadOnlyFile(t *testing.T) {
	outputPath := fmt.Sprintf("%s/output.csv", t.TempDir())

	_, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE, 0444)
	if err != nil {
		t.Fatalf("couldn't create a test config file: %v", err)
	}

	row := []byte("test\n")
	err = appendToOutput(outputPath, row)

	if !errors.Is(err, os.ErrPermission) {
		t.Fatalf("expected path error os.PathError, got %T", err)
	}
}

func TestAppendToOutput(t *testing.T) {
	outputPath := fmt.Sprintf("%s/output.csv", t.TempDir())

	row := []byte("test\n")

	if err := appendToOutput(outputPath, row); err != nil {
		t.Fatalf("couldn't append to output file: %v", err)
	}

	file, err := os.Open(outputPath)
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

func TestGetTimeNoTimezone(t *testing.T) {
	_, err := getTime("aaaaa")
	expected := "unknown time zone aaaaa"

	if err.Error() != expected {
		t.Errorf("expected error '%s', got '%s'", expected, err)
	}
}

func TestGetTime(t *testing.T) {
	tests := map[string]string{
		"Europe/London": "GMT",
		"Europe/Warsaw": "CET",
		"Europe/Kyiv":   "EET",
	}

	for timezone, expectedName := range tests {
		time, _ := getTime(timezone)
		name, _ := time.Zone()

		if name != expectedName {
			t.Errorf("invalid timezone, expected %s/%s, got: %s", timezone, expectedName, name)
		}
	}
}

func TestGetClubListInvalidURL(t *testing.T) {
	_, err := getClubList("::invalid-url", "cookie=abc")

	var urlErr *url.Error
	if !errors.As(err, &urlErr) {
		t.Fatalf("expected error of type *url.Error, got %T", err)
	}
}

func TestGetClubListRequestError(t *testing.T) {
	_, err := getClubList("http://localhost:1", "cookie=abc")

	var urlErr *url.Error
	if !errors.As(err, &urlErr) {
		t.Fatalf("expected error of type *url.Error, got %T", err)
	}
}

func TestGetClubListInvalidResponseStatus(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusForbidden)
			},
		),
	)

	defer ts.Close()

	_, err := getClubList(ts.URL+"/api", "cookie=abc")

	if err.Error() != "response status != 200: 403, body: " {
		t.Fatalf("invalid error, got: %v", err)
	}
}

func TestGetClubListInvalidJSON(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("abc"))
			},
		),
	)

	defer ts.Close()

	_, err := getClubList(ts.URL+"/api", "cookie=abc")

	if err.Error() != "failed to parse json: invalid character 'a' looking for beginning of value" {
		t.Fatalf("invalid error, got: %v", err)
	}
}

func TestGetClubList(t *testing.T) {
	mockResponse := `{"UsersInClubList":[{"ClubAddress":"123 Gym St","UsersCountCurrentlyInClub":10}]}`
	expectedClubList := ClubList{
		Clubs: []ClubData{
			{Address: "123 Gym St", UsersCount: 10},
		},
	}
	expectedCookie := "test_cookie=cookie"

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api" {
					t.Errorf("invalid request path, expected /api, got: %s", r.URL.Path)
				}

				if cookie := r.Header.Get("Cookie"); cookie != expectedCookie {
					t.Errorf("invalid cookie, expected %s, got: %s", expectedCookie, cookie)
				}

				w.Write([]byte(mockResponse))
			},
		),
	)

	defer ts.Close()

	clubs, _ := getClubList(ts.URL+"/api", expectedCookie)

	if !cmp.Equal(clubs, expectedClubList) {
		t.Errorf("invalid result, expected %v, got: %v", expectedClubList, clubs)
	}
}
