package hamtraffic

import (
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

const (
	Freerun         = true
	MaxStationCount = 65536
)

type Config struct {
	Freerun      bool
	StationCount int
}

func NewConfig() *Config {
	var config Config

	fr := os.Getenv("FREERUN")
	if fr == "" {
		config.Freerun = Freerun
	} else {
		f, err := strconv.ParseBool(fr)
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}
		config.Freerun = f
	}

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

	log.Info().Any("config", config).Msg("Configuration complete")

	return &config
}
