package main

import (
	"github.com/andrewyazura/gym-stats/config"
	"github.com/andrewyazura/gym-stats/gymportal"
	"github.com/andrewyazura/gym-stats/storage"
	"log"
)

func main() {
	config, err := config.Load("./config.json")
	if err != nil {
		log.Fatalln(err)
	}

	clubs, err := gymportal.GetClubList(
		config.GymPortal.URL,
		config.GymPortal.Cookies,
	)
	if err != nil {
		log.Fatalln(err)
	}

	storage.WriteAllClubs(config.InfluxDB, clubs)
}
