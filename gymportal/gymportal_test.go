package gymportal

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestClubToPoint(t *testing.T) {
	club := Club{
		Name:      "Gym",
		Address:   "123 Gym St",
		Occupancy: 10,
	}
	point := club.ToPoint()

	tags := point.TagList()
	if len(tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(tags))
	}

	fields := point.FieldList()
	if len(fields) != 1 {
		t.Errorf("expected 1 field, got %d", len(fields))
	}

	if fields[0].Key != "occupancy" {
		t.Errorf("expected field key 'occupancy', got '%s'", fields[0].Key)
	}

	if fields[0].Value != int64(10) {
		t.Errorf("expected field value 10, got %d", fields[0].Value)
	}
}

func TestGetClubListInvalidURL(t *testing.T) {
	_, err := GetClubList("::invalid-url", "cookie=abc")

	var urlErr *url.Error
	if !errors.As(err, &urlErr) {
		t.Fatalf("expected error of type *url.Error, got %T", err)
	}
}

func TestGetClubListRequestError(t *testing.T) {
	_, err := GetClubList("http://localhost:1", "cookie=abc")

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

	_, err := GetClubList(ts.URL+"/api", "cookie=abc")

	if err.Error() != "response status != 200: 403, body: " {
		t.Fatalf("invalid error, got: '%v'", err)
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

	_, err := GetClubList(ts.URL+"/api", "cookie=abc")

	expectedErrorMessage := "failed to parse json: invalid character 'a' looking for beginning of value"
	if err.Error() != expectedErrorMessage {
		t.Fatalf("invalid error, got: '%v'", err)
	}
}

func TestGetClubList(t *testing.T) {
	mockResponse := `{"UsersInClubList":[{"ClubName": "Gym", "ClubAddress":"123 Gym St","UsersCountCurrentlyInClub":10}]}`

	expectedClub := Club{Name: "Gym", Address: "123 Gym St", Occupancy: 10}
	expectedCookie := "test_cookie=cookie"

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api" {
					t.Errorf("expected request path '/api', got: '%s'", r.URL.Path)
				}

				if cookie := r.Header.Get("Cookie"); cookie != expectedCookie {
					t.Errorf("expected cookie '%s', got: '%s'", expectedCookie, cookie)
				}

				w.Write([]byte(mockResponse))
			},
		),
	)

	defer ts.Close()

	clubs, _ := GetClubList(ts.URL+"/api", expectedCookie)
	fmt.Printf("%v\n", clubs)

	if len(clubs) != 1 {
		t.Fatalf("expected club list to have length 1, got length %d", len(clubs))
	}

	if clubs[0] != expectedClub {
		t.Fatalf("expected club '%v', got: '%v'", expectedClub, clubs)
	}
}
