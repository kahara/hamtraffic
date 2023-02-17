package hamtraffic

import "github.com/rs/zerolog/log"

const (
	StationCount = 10000
)

var (
	world    *World
	stations []*Station
)

func Init() {
	world = NewWorld()

	for i := 0; i < StationCount; i++ {
		stations = append(stations, NewStation(world))
	}

	for _, station := range stations {
		log.Info().Str("callsign", station.Callsing).Any("locale", station.Locale).Msg("")
	}
}
