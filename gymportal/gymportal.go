package gymportal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Response struct {
	Clubs ClubList `json:"UsersInClubList"`
}

type ClubList []Club
type Club struct {
	Name      string `json:"ClubName"`
	Address   string `json:"ClubAddress"`
	Occupancy int    `json:"UsersCountCurrentlyInClub"`
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
