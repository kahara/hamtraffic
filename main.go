package hamtraffic

import "github.com/rs/zerolog/log"

var (
	config   *Config
	world    *World
	stations []*Station
)

func Init() {
	config = NewConfig()
	world = NewWorld()

	for i := 0; i < config.StationCount; i++ {
		stations = append(stations, NewStation(config, world))
	}

	for _, station := range stations {
		log.Info().Any("station", station).Msg("")
	}
}
