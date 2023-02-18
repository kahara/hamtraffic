package hamtraffic

import (
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

const (
	MaxStationCount = 65536
)

type Config struct {
	StationCount int
}

func NewConfig() *Config {
	var config Config

	sc := os.Getenv("STATION_COUNT")
	if sc == "" {
		config.StationCount = MaxStationCount
	} else {
		c, err := strconv.Atoi(sc)
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}
		config.StationCount = c
	}

	return &config
}
