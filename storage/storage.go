package storage

import (
	"log"

	"github.com/andrewyazura/gym-occupancy/config"
	"github.com/andrewyazura/gym-occupancy/gymportal"
	"github.com/influxdata/influxdb-client-go/v2"
)

func WriteAllClubs(config config.InfluxDBConfig, clubs gymportal.ClubList) {
	client := influxdb2.NewClient(config.URL, config.AuthToken)
	defer client.Close()

	writeAPI := client.WriteAPI(config.Org, config.Bucket)
	defer writeAPI.Flush()

  errorsCh := writeAPI.Errors()
  go func() {
    for err := range errorsCh {
      log.Printf("influxdb error: %v\n", err)
    }
  }()

	for _, club := range clubs {
		writeAPI.WritePoint(club.ToPoint())
	}
}
