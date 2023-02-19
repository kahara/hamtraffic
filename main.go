package hamtraffic

import (
	"github.com/rs/zerolog/log"
	"time"
)

func Init() {
	var (
		stations        []*Station
		start, deadline time.Time
	)

	config := NewConfig()
	world := NewWorld()

	for i := 0; i < config.StationCount; i++ {
		stations = append(stations, NewStation(config, world))
	}

	for _, station := range stations {
		log.Info().Any("station", station).Msg("")
	}

	// Start running at the start of next minute
	start = time.Now().UTC()
	start = start.Truncate(time.Duration(time.Minute))
	start = start.Add(time.Minute)

	if config.Freerun {
		Run(&start, nil, stations)
	} else {
		deadline = start.Add(*config.Runtime)
		Run(&start, &deadline, stations)
	}
}
