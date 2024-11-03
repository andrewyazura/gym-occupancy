package gymportal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/influxdata/influxdb-client-go/v2"
	write "github.com/influxdata/influxdb-client-go/v2/api/write"
)

type Response struct {
	Clubs ClubList `json:"UsersInClubList"`
}

type ClubList []Club
type Club struct {
	Address   string `json:"ClubAddress"`
	Name      string `json:"ClubName"`
	Occupancy int    `json:"UsersCountCurrentlyInClub"`
}

func (club *Club) ToPoint() *write.Point {
	return influxdb2.NewPoint(
		"club",
		map[string]string{
			"address": club.Address,
			"name":    club.Name,
		},
		map[string]interface{}{
			"occupancy": club.Occupancy,
		},
		time.Now(),
	)
}

func GetClubList(url string, cookies string) (ClubList, error) {
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

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return ClubList{}, fmt.Errorf("failed to parse json: %w", err)
	}

	return response.Clubs, nil
}
