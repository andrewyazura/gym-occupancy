package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Config struct {
	OutputFile string
	Keywords   []string
	URL        string
	Cookies    string
	Timezone   string
}

type ClubList struct {
	Clubs []ClubData `json:"UsersInClubList"`
}

type ClubData struct {
	Address    string `json:"ClubAddress"`
	UsersCount int    `json:"UsersCountCurrentlyInClub"`
}

func loadConfig(path string) (Config, error) {
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

func appendToOutput(path string, row []byte) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	defer file.Close()

	if _, err := file.Write(row); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}

	return nil
}

func getTime(timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}

	return time.Now().In(loc), nil
}

func getClubList(url string, cookies string) (ClubList, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ClubList{}, fmt.Errorf("error creating a request to %s: %w", url, err)
	}

	req.Header.Add("Cookie", cookies)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ClubList{}, fmt.Errorf("error making request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return ClubList{}, fmt.Errorf("response status != 200: %d, body: %s", resp.StatusCode, string(body))
	}

	var list ClubList
  if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return ClubList{}, fmt.Errorf("failed to parse json: %w", err)
	}

	return list, nil
}

func main() {
	var path = flag.String("config", "./config.json", "path to a config.json file")
	flag.Parse()

	config, err := loadConfig(*path)
	if err != nil {
		log.Fatalln(err)
	}

	currentTime, err := getTime(config.Timezone)
	if err != nil {
		log.Fatalln(err)
	}

	clubs, err := getClubList(config.URL, config.Cookies)
	if err != nil {
		log.Fatalln(err)
	}

	for _, club := range clubs.Clubs {
		address := strings.ToLower(club.Address)

		for _, keyword := range config.Keywords {
			if strings.Contains(address, keyword) {
				row := fmt.Sprintf("%v,%v,%d\n", currentTime, address, club.UsersCount)

				if appendToOutput(config.OutputFile, []byte(row)) != nil {
					log.Printf("error while writing to output: %v", err)
				}
			}
		}
	}
}
