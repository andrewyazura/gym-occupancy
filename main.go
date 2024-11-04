package main

import (
	"github.com/andrewyazura/gym-occupancy/config"
	"github.com/andrewyazura/gym-occupancy/gymportal"
	"github.com/andrewyazura/gym-occupancy/storage"
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
